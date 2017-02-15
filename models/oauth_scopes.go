package models

// * Author: Antonio Orozco
// * Date: January 26 2017

// OauthScopes Modelo para administrar información a nivel base de datos.
type OauthScopes struct {
	Scope     string `json:"scope" gorm:"column:scope;size:255;primary_key;not null"`
	IsDefault bool   `json:"is_default" gorm:"column:is_default;type:tinyint(1)"`
}
