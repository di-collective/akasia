package api

import (
	"encoding/json"
	"monorepo/internal/dto"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func (rest *REST) FirebaseAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idToken := strings.TrimSpace(r.URL.Query().Get("idToken"))
	if idToken == "" {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "ID Token should not be empty"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fbaToken, err := rest.oauthVerifier.ValidateIDToken(ctx, idToken)
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Invalid ID Token"})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	signedInEmail, ok := fbaToken.Claims["email"].(string)
	if !ok {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "Unknown Email"})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exists, err := rest.userService.IsHandleExists(ctx, signedInEmail)
	if err == nil && !exists {
		_, err = rest.userService.RegisterUser(ctx, &dto.RequestRegisterUser{
			Provider: fbaToken.Firebase.SignInProvider,
			Email:    signedInEmail,
		})
	}
	if err != nil {
		logrus.Errorf("failed to check user handle: %s; err: %s", signedInEmail, err)
	}

	r.ParseForm()
	r.Form.Set("grant_type", "client_credentials")
	r.Form.Set("client_id", signedInEmail)
	r.Form.Set("client_secret", rest.env.JWTSecret)
	rest.oauthServer.ClientCredentials(w, r)
}
