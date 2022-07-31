package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"testing"
)

func TestGetBearer(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name: `Empty token`,
		},
		{
			name:  `Non-empty token`,
			token: `test_token`,
			want:  `Bearer test_token`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetBearer(tt.token)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTokenFromBearer(t *testing.T) {
	tests := []struct {
		name     string
		bearer   string
		want     string
		wantErr  bool
		errLabel string
	}{
		{
			name:     `Empty bearer`,
			wantErr:  true,
			errLabel: aerror.UserTokenIncorrect,
		},
		{
			name:     `Incorrect bearer`,
			bearer:   `Bearer`,
			wantErr:  true,
			errLabel: aerror.UserTokenIncorrect,
		},
		{
			name:   `Correct bearer`,
			bearer: `Bearer test_token`,
			want:   `test_token`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetTokenFromBearer(tt.bearer)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Label)
			}
		})
	}
}

func TestGetTokenFromUser(t *testing.T) {
	tests := []struct {
		name      string
		user      models.User
		wantEmpty bool
		wantErr   bool
		errLabel  string
	}{
		{
			name: `Token generation`,
			user: models.User{Login: `test_user`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetTokenFromUser(tt.user)
			assert.Equal(t, tt.wantEmpty, got == ``)
			assert.Equal(t, tt.wantErr, err != nil)

			if err != nil {
				assert.Equal(t, tt.errLabel, err.Label)
			}
		})
	}
}

func TestGetUserFromToken(t *testing.T) {
	token, err := GetTokenFromUser(models.User{Login: `test_user`})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		token    string
		want     models.User
		wantErr  bool
		errLabel string
	}{
		{
			name:     `Empty token`,
			wantErr:  true,
			errLabel: aerror.UserTokenIncorrect,
		},
		{
			name:     `Incorrect token`,
			token:    `test_token`,
			wantErr:  true,
			errLabel: aerror.UserTokenIncorrect,
		},
		{
			name:  `Correct token`,
			token: token,
			want:  models.User{Login: `test_user`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetUserFromToken(tt.token)
			fmt.Println(token, tt.token)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Label)
			}
		})
	}
}

func TestIsTokenValid(t *testing.T) {
	token, err := GetTokenFromUser(models.User{Login: `test_user`})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		token    string
		want     bool
		wantErr  bool
		errLabel string
	}{
		{
			name:     `Empty token`,
			wantErr:  true,
			errLabel: aerror.UserTokenInvalid,
		},
		{
			name:     `Invalid token`,
			token:    `test_token`,
			wantErr:  true,
			errLabel: aerror.UserTokenInvalid,
		},
		{
			name:  `Valid token`,
			token: token,
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := IsTokenValid(tt.token)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Label)
			}
		})
	}
}

func Test_keyFn(t *testing.T) {
	tests := []struct {
		name    string
		token   *jwt.Token
		want    interface{}
		wantErr bool
	}{
		{
			name: `Empty token`,
			want: jwtKey,
		},
		{
			name:  `Non-empty token`,
			token: &jwt.Token{},
			want:  jwtKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := keyFn(tt.token)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
