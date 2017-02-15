package main

// * Author: Antonio Orozco
// * Date: January 26 2017

import (
	"log"
	"net/http"

	"github.com/RangelReale/osin"
)

// Start Inicializa el servidor de oAuth2
func main() {
	storage := NewStore()

	err := storage.InitSchemas()

	if err != nil {
		panic(err)
	}

	cfg := osin.NewServerConfig()
	cfg.AllowGetAccessRequest = true
	cfg.AllowClientSecretInParams = true
	cfg.AccessExpiration = 7200
	cfg.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.CLIENT_CREDENTIALS,
		osin.REFRESH_TOKEN,
	}

	// TestStorage implements the "osin.Storage" interface
	server := osin.NewServer(cfg, storage)

	// Authorization code endpoint
	http.HandleFunc("/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		// defer resp.Close()
		var hint string
		if hint = r.FormValue("hint"); hint == "" {
			hint = r.PostFormValue("hint")
		}

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			ar.Authorized = true
			ar.UserData = hint
			server.FinishAuthorizeRequest(resp, r, ar)
		}
		if resp.IsError && resp.InternalError != nil {
			log.Printf("ERROR: %s\n", resp.InternalError)
		}
		osin.OutputJSON(resp, w, r)
	})

	// Access token endpoint
	http.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		// defer resp.Close()

		var hint string
		if hint = r.FormValue("hint"); hint == "" {
			hint = r.PostFormValue("hint")
		}

		ar := server.HandleAccessRequest(resp, r)

		if ar != nil {
			ar.Authorized = true
			ar.UserData = hint

			server.FinishAccessRequest(resp, r, ar)
		}

		if resp.IsError && resp.InternalError != nil {
			log.Printf("ERROR: %s\n", resp.InternalError)
		}

		osin.OutputJSON(resp, w, r)
	})

	http.HandleFunc("/oauth2/passport", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		var token string
		if token = r.FormValue("access_token"); token == "" {
			token = r.PostFormValue("access_token")
		}

		if token != "" {
			access, err := storage.LoadAccess(token)
			if err != nil {
				resp.SetError("invalid_request", "Token is not valid.")
			} else {
				if access.IsExpired() {
					resp.SetError("invalid_request", "Token has expired.")
				} else {
					resp.Output["is_valid"] = true
					resp.Output["token"] = access
				}
			}

		} else {
			resp.SetError("invalid_request", "Required token is missing from the request.")
		}

		osin.OutputJSON(resp, w, r)
	})

	http.ListenAndServe(":8080", nil)
}
