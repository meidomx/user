package model

import (
	"context"

	"github.com/moetang/webapp-scaffold/frmpg"
)

type SsoAppStatus int16

const (
	SsoAppStatusNormal SsoAppStatus = 0
	SsoAppStatusDelete SsoAppStatus = 1
)

type SsoApp struct {
	SsoAppId     int64        `mx.orm:"sso_app_id"`
	AppId        int64        `mx.orm:"app_id"`
	SsoAppStatus SsoAppStatus `mx.orm:"sso_app_status"`
	RedirectUri  string       `mx.orm:"redirect_uri"`
	CallbackUri  string       `mx.orm:"callback_uri"`
	TimeCreated  int64        `mx.orm:"time_created"`
	TimeUpdated  int64        `mx.orm:"time_updated"`
}

func LoadAllSsoAppByAppId(appId int64) (r []*SsoApp, e error) {

	e = frmpg.QueryMulti(db.GetPostgresPool(), &r, context.Background(),
		"select * from sso_app where app_id = $1 and sso_app_status = $2",
		appId, SsoAppStatusNormal)

	return
}
