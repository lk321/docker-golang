package models

// * Author: Antonio Orozco
// * Date: January 26 2017

import "time"

// OauthAuthorizationCodes Modelo para administrar información a nivel base de datos.
type OauthAuthorizationCodes struct {
	AuthorizationCode string    `json:"authorization_code" gorm:"column:authorization_code;primary_key;size:40;not null"`
	ClientID          string    `json:"client_id" gorm:"column:client_id;size:80;not null"`
	LoginHint         string    `json:"login_hint" gorm:"column:login_hint;size:165"`
	RedirectURI       string    `json:"redirect_uri" gorm:"column:redirect_uri;size:255"`
	ExpiresIn         int32     `json:"expires" gorm:"column:expires"`
	State             string    `json:"state" gorm:"column:state;size:255"`
	Scope             string    `json:"scope" gorm:"column:scope;size:255"`
	CreatedAt         time.Time `json:"created" gorm:"column:created;type:timestamp;not null DEFAULT CURRENT_TIMESTAMP"`
}
