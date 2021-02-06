package main

import (
	"github.com/meidomx/user/model"
	"github.com/meidomx/user/restful"
	"github.com/meidomx/user/shared"

	scaffold "github.com/moetang/webapp-scaffold"

	"github.com/gin-gonic/gin"
)

type UserConfig struct {
	User struct {
		RedisConfig struct {
			Host string `toml:"host"`
		} `toml:"redis"`
	} `toml:"user"`
}

func main() {

	webscaf, err := scaffold.NewFromConfigFile("user.toml")
	if err != nil {
		panic(err)
	}
	uc := new(UserConfig)
	if err := scaffold.ReadCustomConfig("user.toml", uc); err != nil {
		panic(err)
	}

	webscaf.GetGin().Use(gin.Logger())
	webscaf.GetGin().Use(gin.Recovery())

	shared.InitRedis(uc.User.RedisConfig.Host, "")
	model.Init(webscaf)
	restful.InitRestful(webscaf)

	err = webscaf.SyncStart()
	if err != nil {
		panic(err)
	}
}
