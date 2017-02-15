package models

// * Author: Antonio Orozco
// * Date: January 26 2017

// OauthClients Modelo para administrar información a nivel base de datos.
type OauthClients struct {
	ClientID     string `json:"client_id" gorm:"column:client_id;size:80;primary_key;not null"`
	ClientSecret string `json:"client_secret" gorm:"column:client_secret;size:80;not null"`
	ClientName   string `json:"client_name" gorm:"column:client_name;size:120;not null"`
	RedirectURI  string `json:"redirect_uri" gorm:"column:redirect_uri;size:255"`
	GrantTypes   string `json:"grant_types" gorm:"column:grant_types;size:80"`
	Hint         string `json:"hint" gorm:"column:hint;size:165"`
	Scope        string `json:"scope" gorm:"column:scope;size:255"`
}
