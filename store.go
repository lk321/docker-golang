package main

// * Author: Antonio Orozco
// * Date: January 26 2017

import (
	"errors"
	"log"
	"time"

	"./models"
	"github.com/RangelReale/osin"
	_ "github.com/go-sql-driver/mysql" // _ Mysql driver
	"github.com/jinzhu/gorm"
)

// Store Struct para administrar oauth2
type Store struct {
	db *gorm.DB
}

// NewStore Crea nuevo store para controlar oauth2
func NewStore() *Store {
	log.Println("oAuth2 Log => NewStore")

	db, err := gorm.Open("mysql", "root:antonio@tcp(127.0.0.1:3306)/oauth2?charset=utf8&parseTime=True")
	if err != nil {
		// error handling
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}

	db.LogMode(true)
	db.SingularTable(true)

	return &Store{db}
}

// InitSchemas Activar migracion
func (s *Store) InitSchemas() error {
	log.Println("oAuth2 Log => InitSchemas")

	return s.db.AutoMigrate(&models.OauthAccessTokens{}, &models.OauthAuthorizationCodes{}, &models.OauthClients{}, &models.OauthRefreshTokens{}, &models.OauthScopes{}).Error
}

// Clone Clona el Store
func (s *Store) Clone() osin.Storage {
	return s
}

// Close Cierra el Store
func (s *Store) Close() {
	s.Close()
}

// GetClient Obtiene cliente de oauth2
func (s *Store) GetClient(id string) (osin.Client, error) {
	log.Printf("oAuth2 Log => GetClient: %s\n", id)
	var c osin.DefaultClient
	client := models.OauthClients{}

	err := s.db.Where("client_id = ?", id).First(&client).Error

	if err != nil {
		return nil, err
	}

	// Mapping client
	c.Id = client.ClientID
	c.Secret = client.ClientSecret
	c.RedirectUri = client.RedirectURI
	c.UserData = client.Hint

	return &c, nil
}

// SaveAuthorize Guarda el codigo de autorizacion de oauth2
func (s *Store) SaveAuthorize(data *osin.AuthorizeData) error {
	log.Printf("oAuth2 Log => SaveAuthorize: %s\n", data.Code)

	var hint string
	if hint = data.UserData.(string); hint == "" {
		hint = data.Client.GetUserData().(string)
	}

	c := models.OauthAuthorizationCodes{
		AuthorizationCode: data.Code,
		ClientID:          data.Client.GetId(),
		LoginHint:         hint,
		RedirectURI:       data.RedirectUri,
		ExpiresIn:         data.ExpiresIn,
		Scope:             data.Scope,
		State:             data.State,
		CreatedAt:         data.CreatedAt,
	}

	return s.db.Create(&c).Error
}

// LoadAuthorize Obtiene el codigo de autorizacion
func (s *Store) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	log.Printf("oAuth2 Log => LoadAuthorize: %s\n", code)
	if code == "" {
		return nil, nil
	}

	var data osin.AuthorizeData
	r := models.OauthAuthorizationCodes{}

	err := s.db.Where("authorization_code = ?", code).First(&r).Error

	data.Code = r.AuthorizationCode
	data.ExpiresIn = r.ExpiresIn
	data.Scope = r.Scope
	data.RedirectUri = r.RedirectURI
	data.State = r.State
	data.CreatedAt = r.CreatedAt
	data.UserData = r.LoginHint

	c, err := s.GetClient(r.ClientID)
	if err != nil {
		return nil, err
	}

	if data.ExpireAt().Before(time.Now()) {
		return nil, errors.New("Token expired at " + data.ExpireAt().String())
	}

	data.Client = c
	return &data, nil
}

// RemoveAuthorize Elimina el codigo de autorizacion
func (s *Store) RemoveAuthorize(code string) error {
	log.Printf("oAuth2 Log => RemoveAuthorize: %s\n", code)

	return s.db.Where("authorization_code = ?", code).Delete(models.OauthAuthorizationCodes{}).Error
}

// SaveAccess Guarda el access token
func (s *Store) SaveAccess(data *osin.AccessData) error {
	log.Printf("oAuth2 Log => SaveAccess: %s\n", data.AccessToken)

	prev := ""
	authorizeData := &osin.AuthorizeData{}

	var hint string
	if hint = data.UserData.(string); hint == "" {
		hint = data.Client.GetUserData().(string)
	}

	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}

	if data.RefreshToken != "" {
		if err := s.SaveRefresh(data.RefreshToken, data.AccessToken, data.Client.GetId(), hint, data.ExpiresIn); err != nil {
			return err
		}
	}

	if data.Client == nil {
		return errors.New("data.Client must not be nil")
	}

	accesToken := models.OauthAccessTokens{
		AccessToken:       data.AccessToken,
		Previous:          prev,
		AuthorizationCode: authorizeData.Code,
		RefreshToken:      data.RefreshToken,
		ClientID:          data.Client.GetId(),
		LoginHint:         hint,
		RedirectURI:       data.RedirectUri,
		ExpiresIn:         data.ExpiresIn,
		Scope:             data.Scope,
		CreatedAt:         data.CreatedAt,
	}

	return s.db.Create(&accesToken).Error
}

// LoadAccess Carga el access token
func (s *Store) LoadAccess(code string) (*osin.AccessData, error) {
	log.Printf("oAuth2 Log => LoadAccess: %s\n", code)
	if code == "" {
		return nil, nil
	}

	result := osin.AccessData{}

	accesToken := models.OauthAccessTokens{}

	err := s.db.Where("access_token = ?", code).First(&accesToken).Error

	if err != nil {
		return nil, err
	}

	result.AccessToken = accesToken.AccessToken
	result.RefreshToken = accesToken.RefreshToken
	result.ExpiresIn = accesToken.ExpiresIn
	result.Scope = accesToken.Scope
	result.RedirectUri = accesToken.RedirectURI
	result.CreatedAt = accesToken.CreatedAt

	result.UserData = accesToken.LoginHint
	client, err := s.GetClient(accesToken.ClientID)
	if err != nil {
		return nil, err
	}

	result.Client = client
	result.AuthorizeData, _ = s.LoadAuthorize(accesToken.AuthorizationCode)
	prevAccess, _ := s.LoadAccess(accesToken.Previous)
	result.AccessData = prevAccess

	return &result, nil
}

// RemoveAccess Elimina access token
func (s *Store) RemoveAccess(token string) error {
	log.Printf("oAuth2 Log => RemoveAccess: %s\n", token)

	return s.db.Where("access_token = ?", token).Delete(models.OauthAccessTokens{}).Error
}

// LoadRefresh Obtiene el refresh token
func (s *Store) LoadRefresh(token string) (*osin.AccessData, error) {
	log.Printf("oAuth2 Log => LoadRefresh: %s\n", token)

	refresh := models.OauthRefreshTokens{}

	s.db.Where("refresh_token = ?", token).First(&refresh)

	return s.LoadAccess(refresh.AccessToken)
}

// RemoveRefresh Elimina refresh token
func (s *Store) RemoveRefresh(token string) error {
	log.Printf("oAuth2 Log => RemoveRefresh: %s\n", token)

	return s.db.Where("refresh_token = ?", token).Delete(models.OauthRefreshTokens{}).Error
}

// SaveRefresh Guarda el refresh en la base de datos
func (s *Store) SaveRefresh(refresh, token, client, hint string, expiresIn int32) (err error) {

	refreshToken := models.OauthRefreshTokens{
		RefreshToken: refresh,
		AccessToken:  token,
		ClientID:     client,
		LoginHint:    hint,
		ExpiresIn:    expiresIn * 2,
	}

	return s.db.Create(&refreshToken).Error
}
