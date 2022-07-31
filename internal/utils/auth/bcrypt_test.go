package auth

import (
	"github.com/stretchr/testify/assert"
	"go-loyalty-system/internal/aerror"
	"testing"
)

func TestCompareHashes(t *testing.T) {
	h, err := HashPassword(`test`)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		reqPwd string
		dbPwd  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name:    `Empty password`,
			wantErr: aerror.UserPasswordIncorrect,
		},
		{
			name: `Non-matching passwords`,
			args: args{
				reqPwd: "test1",
				dbPwd:  "test2",
			},
			wantErr: aerror.UserPasswordIncorrect,
		},
		{
			name: `Matching passwords`,
			args: args{
				reqPwd: `test`,
				dbPwd:  h,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := CompareHashes(tt.args.reqPwd, tt.args.dbPwd)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr == ``, err == nil)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Label)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name      string
		pass      string
		wantEmpty bool
		wantErr   string
	}{
		{
			name:      `Empty password`,
			wantEmpty: true,
			wantErr:   aerror.UserPasswordIncorrect,
		},
		{
			name: `Non-empty password`,
			pass: `test`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := HashPassword(tt.pass)
			assert.Equal(t, tt.wantEmpty, got == ``)
			assert.Equal(t, tt.wantErr == ``, err == nil)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Label)
			}
		})
	}
}
