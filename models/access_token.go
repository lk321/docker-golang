package models

// * Author: Antonio Orozco
// * Date: January 26 2017

import "time"

// OauthAccessTokens Modelo para administrar información a nivel base de datos.
type OauthAccessTokens struct {
	AccessToken       string    `json:"access_token" gorm:"column:access_token;primary_key;size:40;not null"`
	Previous          string    `json:"previous" gorm:"column:previous;size:40;"`
	AuthorizationCode string    `json:"authorization_code" gorm:"column:authorization_code;size:40"`
	RefreshToken      string    `json:"refresh_token" gorm:"column:refresh_token;size:40"`
	ClientID          string    `json:"client_id" gorm:"column:client_id;size:80;not null"`
	LoginHint         string    `json:"login_hint" gorm:"column:login_hint;size:165"`
	ExpiresIn         int32     `json:"expires" gorm:"column:expires"`
	RedirectURI       string    `json:"redirect_uri" gorm:"column:redirect_uri;size:255"`
	Scope             string    `json:"scope" gorm:"column:scope;size:255"`
	CreatedAt         time.Time `json:"created" gorm:"column:created;type:timestamp;not null DEFAULT CURRENT_TIMESTAMP"`
}
