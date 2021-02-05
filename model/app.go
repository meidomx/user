package model

import (
	"context"
	"time"

	"github.com/moetang/webapp-scaffold/frmpg"
)

type AppStatus int16
type AppTokenType int16
type AppTokenStatus int16

const (
	AppStatusNormal AppStatus = 0
	AppStatusDelete AppStatus = 1
)

const (
	AppTokenTypeTokenOnly         AppTokenType = 1
	AppTokenTypeTokenWithSecurity AppTokenType = 2
)

const (
	AppTokenStatusNormal AppTokenStatus = 0
	AppTokenStatusDelete AppTokenStatus = 1
)

type AppApp struct {
	AppId       int64     `mx.orm:"app_id"`
	AppName     string    `mx.orm:"app_name"`
	AppStatus   AppStatus `mx.orm:"app_status"`
	TimeCreated int64     `mx.orm:"time_created"`
	TimeUpdated int64     `mx.orm:"time_updated"`
}

type AppToken struct {
	TokenId            int64          `mx.orm:"token_id"`
	AppId              int64          `mx.orm:"app_id"`
	Token              string         `mx.orm:"token"`
	SecurityValue      *string        `mx.orm:"security_value"`
	TokenType          AppTokenType   `mx.orm:"token_type"`
	TokenStatus        AppTokenStatus `mx.orm:"token_status"`
	ExpirydateInMillis *int64         `mx.orm:"expirydate_millis"`
	TimeCreated        int64          `mx.orm:"time_created"`
	TimeUpdated        int64          `mx.orm:"time_updated"`
}

func LoadAppAppByAppName(appName string) (*AppApp, error) {
	var app = new(AppApp)
	err := frmpg.QuerySingle(db.GetPostgresPool(), app, context.Background(),
		"select * from app_app where app_name = $1 and app_status = $2",
		appName, AppStatusNormal)
	if err != nil {
		return nil, err
	} else {
		return app, nil
	}
}

func LoadAppTokenByToken(token string) (*AppToken, error) {
	var t = new(AppToken)
	err := frmpg.QuerySingle(db.GetPostgresPool(), t, context.Background(),
		"select * from app_token where token = $1 and token_status = $2",
		token, AppTokenStatusNormal)
	if err != nil {
		return nil, err
	} else {
		return t, nil
	}
}

func LoadAppAppByAppId(appId int64) (*AppApp, error) {
	var app = new(AppApp)
	err := frmpg.QuerySingle(db.GetPostgresPool(), app, context.Background(),
		"select * from app_app where app_id = $1 and app_status = $2",
		appId, AppStatusNormal)
	if err != nil {
		return nil, err
	} else {
		return app, nil
	}
}

func ValidateAppNameAndToken(appName, token string) (bool, error) {
	var t = new(AppToken)
	err := frmpg.QuerySingle(db.GetPostgresPool(), t, context.Background(),
		"select * from app_token where token = $1 and token_status = $2 and token_type = $3",
		token, AppTokenStatusNormal, AppTokenTypeTokenOnly)
	if err == frmpg.ErrNoRecordFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// validate time
	if t.ExpirydateInMillis != nil {
		now := time.Now()
		if *t.ExpirydateInMillis <= now.Unix()*1000+now.UnixNano()/1000000 {
			return false, nil
		}
	}

	a, err := LoadAppAppByAppId(t.AppId)
	if err == frmpg.ErrNoRecordFound {
		return false, nil
	}
	if a != nil {
		return true, nil
	} else {
		return false, nil
	}
}
