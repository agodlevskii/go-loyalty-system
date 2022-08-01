package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"regexp"
	"testing"
)

var user = models.User{Login: `test`, Password: `test`}

func TestDBUser_Add(t *testing.T) {
	tests := []struct {
		name    string
		stored  models.User
		user    models.User
		wantErr string
	}{
		{
			name: `Add new user`,
			user: user,
		},
		{
			name:    `Add existing user`,
			user:    user,
			stored:  user,
			wantErr: aerror.UserAdd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initUserRepo(t, tt.stored)
			defer r.db.Close()

			addUserExpect(mock, tt.user, true)
			err := r.Add(tt.user)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBUser_Find(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		stored  models.User
		want    models.User
		wantErr string
	}{
		{
			name:    `Missing user`,
			user:    user.Login,
			wantErr: aerror.UserFind,
		},
		{
			name:   `Existing user`,
			stored: user,
			want:   user,
			user:   user.Login,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initUserRepo(t, tt.stored)
			defer r.db.Close()

			eq := mock.ExpectQuery(regexp.QuoteMeta(UserFind)).WithArgs(tt.user)
			if tt.stored.Login != `` {
				eq.WillReturnRows(sqlmock.NewRows([]string{`password`}).AddRow(tt.stored.Password))
			} else {
				eq.WillReturnError(aerror.NewError(aerror.UserFind, sql.ErrNoRows))
			}

			got, err := r.Find(tt.user)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestNewDBUserStorage(t *testing.T) {
	db, mock := getMock(t)
	tests := []struct {
		name    string
		db      *sql.DB
		want    DBUser
		wantErr string
	}{
		{
			name: `Create storage`,
			db:   db,
			want: DBUser{db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock.ExpectExec(regexp.QuoteMeta(UserTableCreate)).WillReturnResult(sqlmock.NewResult(1, 1))
			got, err := NewDBUserStorage(tt.db)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func initUserRepo(t *testing.T, init models.User) (DBUser, sqlmock.Sqlmock) {
	db, mock := getMock(t)
	r := DBUser{db}
	if init.Login != `` {
		addUserExpect(mock, init, false)
		if err := r.Add(init); err != nil {
			t.Fatal(err)
		}
	}
	return r, mock
}

func addUserExpect(mock sqlmock.Sqlmock, user models.User, duplicate bool) {
	eq := mock.ExpectQuery(regexp.QuoteMeta(UserAdd)).
		WithArgs(user.Login, user.Password)

	if duplicate {
		eq.WillReturnError(aerror.NewError(aerror.UserAdd, sql.ErrNoRows))
	} else {
		eq.WillReturnRows(sqlmock.NewRows([]string{`name`}).AddRow(user.Login))
	}
}
