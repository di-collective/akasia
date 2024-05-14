package api

import (
	"encoding/json"
	"io"
	"monorepo/internal/dto"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func (rest *REST) FirebaseAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
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
	if err != nil {
		json.NewEncoder(w).Encode(dto.Object[any]{Error: "User Not Exist"})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !exists {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Object[any]{Error: "Failed to Parse Payload"})
			return
		}

		var reqRegister dto.RequestRegisterUser
		json.Unmarshal(payload, &reqRegister)

		reqRegister.Provider = fbaToken.Firebase.SignInProvider
		reqRegister.Email = signedInEmail

		_, err = rest.userService.RegisterUser(ctx, &reqRegister)
		if err != nil {
			logrus.Errorf("failed to register user: %s; err: %s", reqRegister.Email, err)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Object[any]{Error: "Failed to Register User"})
			return
		}

	}

	// login
	r.ParseForm()
	r.Form.Set("grant_type", "client_credentials")
	r.Form.Set("client_id", signedInEmail)
	r.Form.Set("client_secret", rest.env.JWTSecret)
	rest.oauthServer.ClientCredentials(w, r)
}
