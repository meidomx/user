package model

import (
	"context"
	"errors"
	"time"

	"github.com/moetang/webapp-scaffold/frmpg"

	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/jackc/pgx/v4"
)

type UserType int16
type UserStatus int16
type UserCredentialType int16
type UserCredentialStatus int16

const (
	UserTypeStandard UserType = 1
)

const (
	UserStatusNormal  UserStatus = 0
	UserStatusDeleted UserStatus = 1
)

const (
	UserCredentialTypeUserNamePassword UserCredentialType = 1
)

const (
	UserCredentialStatusNormal  UserCredentialStatus = 0
	UserCredentialStatusDeleted UserCredentialStatus = 1
)

type UserUser struct {
	UserId      int64      `mx.orm:"user_id"`
	UserType    UserType   `mx.orm:"user_type"`
	UserStatus  UserStatus `mx.orm:"user_status"`
	UserTag1    int64      `mx.orm:"user_tag1"`
	UserTag2    int64      `mx.orm:"user_tag2"`
	UserName    string     `mx.orm:"user_name"`
	DisplayName string     `mx.orm:"display_name"`
	TimeCreated int64      `mx.orm:"time_created"`
	TimeUpdated int64      `mx.orm:"time_updated"`
}

type UserCredential struct {
	CredentialId     int64                `mx.orm:"credential_id"`
	UserId           int64                `mx.orm:"user_id"`
	CredentialType   UserCredentialType   `mx.orm:"credential_type"`
	CredentialKey    string               `mx.orm:"credential_key"`
	CredentialValue  *string              `mx.orm:"credential_value"`
	CredentialStatus UserCredentialStatus `mx.orm:"credential_status"`
	TimeCreated      int64                `mx.orm:"time_created"`
	TimeUpdated      int64                `mx.orm:"time_updated"`
}

func CreateUser(username, password, displayName string) (u *UserUser, e error) {
	hash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		return nil, err
	}
	encPw := string(hash)

	e = frmpg.DoTxDefault(db.GetPostgresPool(), func(tx pgx.Tx) error {
		now := time.Now()
		rs, err := tx.Query(context.Background(),
			"insert into user_user(user_type, user_status, user_tag1, user_tag2, time_created, time_updated, user_name, display_name) values ($1, $2, $3, $4, $5, $6, $7, $8) returning user_id",
			UserTypeStandard, UserStatusNormal, 0, 0, UnixTime(now), UnixTime(now), username, displayName)
		if err != nil {
			return err
		}
		defer rs.Close()

		var userId int64
		if rs.Next() {
			if err := rs.Scan(&userId); err != nil {
				return err
			}
		} else {
			return ErrNoRecordFound
		}
		rs.Close()

		_, err = tx.Exec(context.Background(),
			"insert into user_credential(user_id, credential_type, credential_key, credential_value, credential_status, time_created, time_updated) values($1, $2, $3, $4, $5, $6, $7)",
			userId, UserCredentialTypeUserNamePassword, username, encPw, UserCredentialStatusNormal, UnixTime(now), UnixTime(now))
		if err != nil {
			return err
		}

		u = new(UserUser)
		err = frmpg.QuerySingle(tx, u, context.Background(),
			"select * from user_user where user_id = $1",
			userId)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func QueryUserById(userId int64) (*UserUser, error) {
	u := new(UserUser)
	err := frmpg.QuerySingle(db.GetPostgresPool(), u, context.Background(),
		"select * from user_user where user_id = $1 and user_status = $2",
		userId, UserStatusNormal)
	if err == ErrNoRecordFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func VerifyUserNameAndPasswordWithUserReturned(username, password string) (u *UserUser, e error) {
	c := new(UserCredential)

	err := frmpg.DoTxDefault(db.GetPostgresPool(), func(tx pgx.Tx) error {

		err := frmpg.QuerySingle(tx, c, context.Background(),
			"select * from user_credential where credential_key = $1 and credential_status = $2 and credential_type = $3",
			username, UserCredentialStatusNormal, UserCredentialTypeUserNamePassword)
		if err != nil {
			return err
		}

		u = new(UserUser)
		err = frmpg.QuerySingle(tx, u, context.Background(),
			"select * from user_user where user_id = $1 and user_status = $2",
			c.UserId, UserStatusNormal)
		if err != nil {
			return err
		}

		return nil
	})
	if err == ErrNoRecordFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	err = scrypt.CompareHashAndPassword([]byte(*c.CredentialValue), []byte(password))
	if err == scrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}
	return
}
