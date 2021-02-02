package model

import (
	scaffold "github.com/moetang/webapp-scaffold"
)

var UserServToken = ""

var db scaffold.PostgresApi

func Init(webscaf *scaffold.WebappScaffold) {
	db = webscaf
}
