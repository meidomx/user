package main

import (
	"github.com/meidomx/user/model"
	"github.com/meidomx/user/restful"

	scaffold "github.com/moetang/webapp-scaffold"

	"github.com/gin-gonic/gin"
)

func main() {

	webscaf, err := scaffold.NewFromConfigFile("user.toml")
	if err != nil {
		panic(err)
	}

	var _ = webscaf

	webscaf.GetGin().Use(gin.Logger())
	webscaf.GetGin().Use(gin.Recovery())

	model.Init(webscaf)
	restful.InitRestful(webscaf)

	err = webscaf.SyncStart()
	if err != nil {
		panic(err)
	}
}
