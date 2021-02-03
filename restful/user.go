package restful

import (
	"log"
	"net/http"
	"strconv"

	"github.com/meidomx/user/model"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	UserType    int16  `json:"user_type" binding:"required"`
	UserName    string `json:"user_name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type UserResult struct {
	UserId      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	DisplayName string `json:"display_name"`
}

func CreateUser(ctx *gin.Context) {
	var u = new(CreateUserRequest)
	if err := ctx.ShouldBindJSON(u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorcode":    "0001",
			"errormessage": "parameter invalid",
		})
		return
	}

	switch u.UserType {
	case 1:
		// username/password
		uu, err := model.CreateUser(u.UserName, u.Password, u.DisplayName)
		if err != nil {
			log.Println("[ERROR] create failed.", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errorcode":    "0003",
				"errormessage": "create user failed",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user_id": uu.UserId,
		})
		return
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorcode":    "0002",
			"errormessage": "user type unsupported",
		})
		return
	}

}

func GetUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		log.Println("[ERROR] parse request failed.", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorcode":    "0001",
			"errormessage": "parameter invalid",
		})
		return
	}

	u, err := model.QueryUserById(userId)
	if err != nil {
		log.Println("[ERROR] query user by id failed.", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorcode":    "0004",
			"errormessage": "query user failed",
		})
		return
	}

	if u == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"errorcode":    "0005",
			"errormessage": "no user found",
		})
		return
	}

	ctx.JSON(http.StatusOK, UserResult{
		UserId:      u.UserId,
		UserName:    u.UserName,
		DisplayName: u.DisplayName,
	})
}
