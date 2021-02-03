package restful

import (
	"net/http"

	"github.com/meidomx/user/model"

	"github.com/gin-gonic/gin"
)

type AuthPasswordRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AuthPassword(ctx *gin.Context) {
	var u = new(AuthPasswordRequest)
	if err := ctx.ShouldBindJSON(u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errorcode":    "0001",
			"errormessage": "parameter invalid",
		})
		return
	}

	uu, err := model.VerifyUserNameAndPasswordWithUserReturned(u.UserName, u.Password)
	if err != nil || uu == nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"errorcode":    "0006",
			"errormessage": "not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id": uu.UserId,
	})
}
