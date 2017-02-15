package models

import "time"

// * Author: Antonio Orozco
// * Date: January 26 2017

// OauthRefreshTokens Modelo para administrar información a nivel base de datos.
type OauthRefreshTokens struct {
	RefreshToken string    `json:"refresh_token" gorm:"column:refresh_token;primary_key;size:40;not null"`
	AccessToken  string    `json:"access_token" gorm:"column:access_token;size:40;not null"`
	ClientID     string    `json:"client_id" gorm:"column:client_id;size:80;not null"`
	LoginHint    string    `json:"login_hint" gorm:"column:login_hint;size:165"`
	ExpiresIn    int32     `json:"expires" gorm:"column:expires"`
	CreatedAt    time.Time `json:"created" gorm:"column:created;type:timestamp;not null DEFAULT CURRENT_TIMESTAMP"`
}
