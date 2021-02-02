package restful

import (
	"github.com/meidomx/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initApi(g *gin.Engine) {
	grp := g.Group("/pub_api/v1")

	grp.Use(func(context *gin.Context) {
		token := context.Request.Header.Get("Token")

		if len(token) == 0 {
			context.JSON(http.StatusForbidden, gin.H{})
			return
		}

		ok, err := model.ValidateAppNameAndToken("user", token)
		if err != nil {
			panic(err)
		}

		if ok {
			context.Next()
			return
		} else {
			context.JSON(http.StatusForbidden, gin.H{})
			return
		}
	})

	// create user
	grp.POST("/user")
	// get user info
	grp.GET("/user/:user_id")

	// auth using password method
	grp.POST("/auth/password")
}
