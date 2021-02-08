package restful

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/meidomx/user/model"
	"github.com/meidomx/user/shared"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
)

func initSso(g *gin.Engine) {

	grp := g.Group("/sso/v1")

	grp.GET("/auth", Auth)
	grp.POST("/login", Login)
	grp.POST("/token", Token)
	grp.GET("/user_info", SsoUserInfo)
}

type CacheItem struct {
	ClientId     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectUri  string   `json:"redirect_uri"`
	Scope        []string `json:"scope"`
	State        string   `json:"state"`
	AuthCode     string   `json:"auth_code"`
	CallbackUrl  string   `json:"callback_url"`
	SsoAppId     int64    `json:"sso_app_id"`
	AppTokenId   int64    `json:"app_token_id"`

	UserId          int64  `json:"user_id"`
	UserName        string `json:"user_name"`
	UserDisplayName string `json:"user_display_name"`
}

type UserInfo struct {
	UserId          int64  `json:"user_id"`
	UserName        string `json:"user_name"`
	UserDisplayName string `json:"user_display_name"`
}

type AuthRequest struct {
	ResponseType string   `form:"response_type" binding:"required" json:"response_type"`
	ClientId     string   `form:"client_id" binding:"required" json:"client_id"`
	RedirectUri  string   `form:"redirect_uri" binding:"required" json:"redirect_uri"`
	Scope        []string `form:"scope" binding:"required" json:"scope"`
	State        string   `form:"state" binding:"required" json:"state"`
}

type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code" binding:"required"`
	RedirectUri  string `form:"redirect_uri" binding:"required"`
	ClientId     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
}

type TokenReply struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"` // second
}

type TokenFailReply struct {
	Error string `json:"error"`
}

type GeneralFailReply struct {
	ErrorMessage string
}

func Auth(ctx *gin.Context) {
	ar := new(AuthRequest)
	if err := ctx.ShouldBindQuery(ar); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "param invalid",
		})
		return
	}

	if ar.ResponseType != "code" {
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "param not supported",
		})
		return
	}

	//TODO scope support
	var _ = ar.Scope

	//TODO check token type
	//TODO should confirm it is used for sso
	at, err := model.LoadAppTokenByToken(ar.ClientId)
	if err != nil {
		log.Println("[ERROR] load app token by token failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}
	aa, err := model.LoadAppAppByAppId(at.AppId)
	if err != nil {
		log.Println("[ERROR] load AppApp by appId failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}

	//TODO check redirectUrl
	//TODO currently only support wildcard '*'
	r, err := model.LoadAllSsoAppByAppId(aa.AppId)
	if err != nil {
		log.Println("[ERROR] load all SsoApp by appId failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}
	var found *model.SsoApp
	for _, i := range r {
		if i.RedirectUri == "*" {
			found = i
			break
		}
	}
	if found == nil {
		log.Println("[ERROR] no redirect url config found.")
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "param failed",
		})
		return
	}

	ci := new(CacheItem)
	ci.ClientId = ar.ClientId
	//ci.ClientSecret = *at.SecurityValue
	ci.RedirectUri = ar.RedirectUri
	ci.Scope = ar.Scope
	ci.State = ar.State
	ci.CallbackUrl = found.CallbackUri
	ci.SsoAppId = found.SsoAppId
	ci.AppTokenId = at.TokenId

	data, err := json.Marshal(ci)
	if err != nil {
		log.Println("[ERROR] marshal redis obj failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}

	uu, err := uuid.NewV4()
	if err != nil {
		log.Println("[ERROR] new uuid failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}
	c := shared.RedisClient.Set(context.Background(), makeKey("authreq", aa.AppId, uu.String()), data, 10*time.Minute)
	_, err = c.Result()
	if err != nil {
		log.Println("[ERROR] store redis failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}

	ctx.HTML(http.StatusOK, "sso_auth.html", gin.H{
		"ItemId": makeKey("authreq", aa.AppId, uu.String()),
		"AppId":  aa.AppId,
	})
}

func makeKey(t string, appId int64, timestamp string) string {
	return fmt.Sprint(t, "_", appId, "_", timestamp)
}

func Login(ctx *gin.Context) {
	//TODO need captcha

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	authid := ctx.PostForm("authid")

	if username == "" || password == "" || authid == "" {
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "param invalid",
		})
		return
	}

	r, err := shared.RedisClient.Get(context.Background(), authid).Result()
	if err == redis.Nil {
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "invalid token",
		})
		return
	} else if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}
	// delete authid anyway
	shared.RedisClient.Del(context.Background(), authid)

	var ci = new(CacheItem)
	err = json.Unmarshal([]byte(r), ci)
	if err != nil {
		log.Println("[ERROR] unmarshal redis obj failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}

	useruser, err := model.VerifyUserNameAndPasswordWithUserReturned(username, password)
	if err != nil || useruser == nil {
		log.Println("[ERROR] login failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "login failed",
		})
		return
	}
	ci.UserId = useruser.UserId
	ci.UserName = useruser.UserName
	ci.UserDisplayName = useruser.DisplayName

	uu, err := uuid.NewV4()
	if err != nil {
		log.Println("[ERROR] new uuid failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}
	ci.AuthCode = uu.String()
	data, err := json.Marshal(ci)
	if err != nil {
		log.Println("[ERROR] marshal redis obj failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}
	c := shared.RedisClient.Set(context.Background(), makeKey("authcode", 0, ci.AuthCode), data, 10*time.Minute)
	_, err = c.Result()
	if err != nil {
		log.Println("[ERROR] store redis failed.", err)
		ctx.HTML(http.StatusBadRequest, "error.html", GeneralFailReply{
			ErrorMessage: "internal error",
		})
		return
	}

	ctx.Redirect(http.StatusFound, fmt.Sprint(ci.RedirectUri, "?code=", ci.AuthCode, "&state=", ci.State))
}

func Token(ctx *gin.Context) {
	var to = new(TokenRequest)
	if err := ctx.ShouldBindQuery(to); err != nil {
		if err = ctx.ShouldBind(to); err != nil {
			ctx.JSON(http.StatusBadRequest, TokenFailReply{
				Error: "param invalid",
			})
			return
		}
	}
	if to.GrantType != "authorization_code" {
		log.Println("[ERROR] grant type error.")
		ctx.JSON(http.StatusBadRequest, TokenFailReply{
			Error: "param invalid",
		})
		return
	}

	r, err := shared.RedisClient.Get(context.Background(), makeKey("authcode", 0, to.Code)).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "invalid code",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}
	// delete code anyway
	shared.RedisClient.Del(context.Background(), makeKey("authcode", 0, to.Code))

	var ci = new(CacheItem)
	err = json.Unmarshal([]byte(r), ci)
	if err != nil {
		log.Println("[ERROR] unmarshal redis obj failed.", err)
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	if ci.RedirectUri != to.RedirectUri {
		log.Println("[ERROR] redirect url not match")
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}
	if ci.ClientId != to.ClientId {
		log.Println("[ERROR] client id not match")
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	at, err := model.LoadAppTokenByToken(ci.ClientId)
	if err != nil {
		log.Println("[ERROR] load app token by token failed.", err)
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	if at.SecurityValue == nil || *at.SecurityValue != to.ClientSecret {
		log.Println("[ERROR] security value not match.")
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	//TODO need to persis access token
	uu, err := uuid.NewV4()
	if err != nil {
		log.Println("[ERROR] new uuid failed.", err)
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	now := time.Now()
	oneh := now.Add(1 * time.Hour)

	ui := new(UserInfo)
	ui.UserId = ci.UserId
	ui.UserName = ci.UserName
	ui.UserDisplayName = ci.UserDisplayName

	uis, err := json.Marshal(ui)
	if err != nil {
		log.Println("[ERROR] marshal UserInfo failed.", err)
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}
	cmd := shared.RedisClient.Set(context.Background(), makeUserInfoKey(uu.String()), uis, 1*time.Hour) // 1h ttl
	if _, err := cmd.Result(); err != nil {
		log.Println("[ERROR] set UserInfo failed.", err)
		ctx.JSON(http.StatusInternalServerError, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	tr := new(TokenReply)
	tr.AccessToken = uu.String()
	tr.ExpiresIn = int(oneh.Unix())
	ctx.JSON(http.StatusOK, tr)
}

func makeUserInfoKey(s string) string {
	return "ssouserinfo_" + s
}

func SsoUserInfo(ctx *gin.Context) {
	a := ctx.GetHeader("Authorization")
	if len(a) == 0 {
		ctx.JSON(http.StatusForbidden, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	ss := strings.Split(a, " ")
	if len(ss) != 2 {
		ctx.JSON(http.StatusForbidden, TokenFailReply{
			Error: "internal error",
		})
		return
	}

	cmd := shared.RedisClient.Get(context.Background(), makeUserInfoKey(ss[1]))
	if d, err := cmd.Result(); err != nil {
		ctx.JSON(http.StatusForbidden, TokenFailReply{
			Error: "internal error",
		})
		return
	} else {
		var ui = new(UserInfo)
		err := json.Unmarshal([]byte(d), ui)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TokenFailReply{
				Error: "internal error",
			})
			return
		}
		ctx.JSON(http.StatusOK, ui)
		return
	}

	panic(errors.New("cannot reach here"))
}
