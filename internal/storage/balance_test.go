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

func TestDBBalance_Get(t *testing.T) {
	tb := models.Balance{
		Current:   100,
		Withdrawn: 50,
		User:      "test",
	}
	tests := []struct {
		name    string
		user    string
		stored  models.Balance
		want    models.Balance
		wantErr string
	}{
		{
			name:    "Non-existing entry",
			user:    tb.User,
			wantErr: aerror.BalanceGet,
		},
		{
			name:   "Existing entry",
			user:   tb.User,
			stored: tb,
			want:   tb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initBalanceRepo(t, tt.stored)
			defer r.db.Close()

			eq := mock.ExpectQuery(regexp.QuoteMeta(BalanceGet)).WithArgs(tt.user)
			if tt.stored.User != "" {
				row := mock.NewRows([]string{"user", "current", "withdrawn"}).AddRow(tt.stored.User, tt.stored.Current, tt.stored.Withdrawn)
				eq.WillReturnRows(row)
			}

			got, err := r.Get(tt.user)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBBalance_Set(t *testing.T) {
	tb := models.Balance{
		Current:   100,
		Withdrawn: 50,
		User:      "test",
	}
	tests := []struct {
		name    string
		b       models.Balance
		stored  models.Balance
		wantErr string
	}{
		{
			name: "New balance",
			b:    tb,
		},
		{
			name:   "Existing balance",
			b:      tb,
			stored: tb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initBalanceRepo(t, tt.stored)
			defer r.db.Close()

			setBalanceExpect(mock, tt.b)
			err := r.Set(tt.b)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestNewDBBalanceStorage(t *testing.T) {
	db, mock := getMock(t)
	tests := []struct {
		name    string
		db      *sql.DB
		want    DBBalance
		wantErr string
	}{
		{
			name: "Create storage",
			db:   db,
			want: DBBalance{db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock.ExpectExec(regexp.QuoteMeta(BalanceTableCreate)).WillReturnResult(sqlmock.NewResult(1, 1))
			got, err := NewDBBalanceStorage(tt.db)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func initBalanceRepo(t *testing.T, init models.Balance) (DBBalance, sqlmock.Sqlmock) {
	db, mock := getMock(t)
	r := DBBalance{db}
	if init.User != "" {
		setBalanceExpect(mock, init)
		if err := r.Set(init); err != nil {
			t.Fatal(err)
		}
	}

	return r, mock
}

func setBalanceExpect(mock sqlmock.Sqlmock, b models.Balance) {
	mock.ExpectExec(regexp.QuoteMeta(BalanceSet)).
		WithArgs(b.User, b.Current, b.Withdrawn).
		WillReturnResult(sqlmock.NewResult(1, 1))
}
