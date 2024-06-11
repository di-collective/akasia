package service

import (
	"context"
	"errors"
	"fmt"
	"monorepo/internal/config"
	"monorepo/pkg/common"
	"monorepo/services/user/models"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/go-chi/oauth"
)

type OauthVerifier struct {
	fbaClient *auth.Client
	env       *config.Environment
	tables    struct {
		user common.Repository[models.User, string]
	}
}

func NewOauthVerifier(
	tbUser common.Repository[models.User, string],
	fbaClient *auth.Client,
	env *config.Environment,
) *OauthVerifier {
	verifier := &OauthVerifier{fbaClient: fbaClient, env: env}
	verifier.tables.user = tbUser

	return verifier
}

func (verifier *OauthVerifier) ValidateIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return verifier.fbaClient.VerifyIDToken(ctx, idToken)
}

// ValidateUser validates username and password returning an error if the user credentials are wrong
func (verifier *OauthVerifier) ValidateUser(username, password, scope string, r *http.Request) error {
	user, err := verifier.tables.user.List(context.Background(), &common.FilterOptions{
		Sort:   []exp.OrderedExpression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("handle").Eq(username)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return err
	}

	if len(user) == 0 {
		err = errors.New("user does not exist")
		return err
	}

	// TODO: compare password
	if username == user[0].Handle {
		return nil
	}

	// if username == user[0].Handle && password == "12345" {
	// 	return nil
	// }

	return errors.New("invalid user")
}

// ValidateClient validates clientID and secret returning an error if the client credentials are wrong
func (verifier *OauthVerifier) ValidateClient(clientID, clientSecret, scope string, r *http.Request) error {
	if clientSecret == verifier.env.JWTSecret {
		return nil
	}

	return errors.New("invalid client")
}

// ValidateCode validates token ID
func (*OauthVerifier) ValidateCode(clientID, clientSecret, code, redirectURI string, r *http.Request) (string, error) {
	return "", nil
}

// AddClaims provides additional claims to the token
func (verifier *OauthVerifier) AddClaims(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	ctx := r.Context()
	existing, err := verifier.tables.user.List(ctx, &common.FilterOptions{
		Page: 1, Limit: 1,
		Select: []any{"id", "handle", "created_at"},
		Filter: []exp.Expression{
			goqu.C("handle").Eq(credential),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w; %w", ErrRepositoryQueryFail, err)
	} else if len(existing) <= 0 {
		return nil, ErrNoResult
	}

	user := existing[0]
	claims := map[string]string{
		"x-hasura-default-role": "user",
		"x-hasura-user-id":      user.ID,
	}

	return claims, nil
}

// AddProperties provides additional information to the token response
func (*OauthVerifier) AddProperties(tokenType oauth.TokenType, credential, tokenID, scope string, r *http.Request) (map[string]string, error) {
	props := make(map[string]string)
	return props, nil
}

// ValidateTokenID validates token ID
func (*OauthVerifier) ValidateTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}

// StoreTokenID saves the token id generated for the user
func (*OauthVerifier) StoreTokenID(tokenType oauth.TokenType, credential, tokenID, refreshTokenID string) error {
	return nil
}
