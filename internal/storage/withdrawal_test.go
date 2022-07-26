package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-loyalty-system/internal/aerror"
	"go-loyalty-system/internal/models"
	"regexp"
	"testing"
	"time"
)

type storedData struct {
	b models.Balance
	w models.Withdrawal
}

var tb = models.Balance{
	Current:   100,
	Withdrawn: 0,
	User:      "test",
}

var tw = models.Withdrawal{
	Order:       "test",
	Sum:         0,
	ProcessedAt: time.Now().Format(time.RFC3339),
	User:        "test",
}

func TestDBWithdrawal_Add(t *testing.T) {
	tests := []struct {
		name    string
		stored  storedData
		w       models.Withdrawal
		wantErr string
	}{
		{
			name: "Existing withdrawal",
			stored: storedData{
				b: tb,
				w: tw,
			},
			w:       tw,
			wantErr: aerror.WithdrawalAdd,
		}, {
			name: "Non-existing withdrawal",
			w:    tw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initWithdrawalRepo(t, tt.stored)
			defer r.db.Close()
			ed := storedData{
				b: tt.stored.b,
				w: tt.w,
			}
			expectWithdrawalAdd(mock, ed, tt.stored.w.User != "")
			err := r.Add(tt.w)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBWithdrawal_Find(t *testing.T) {
	tests := []struct {
		name    string
		order   string
		stored  storedData
		want    models.Withdrawal
		wantErr string
	}{
		{
			name: "Existing withdrawal",
			stored: storedData{
				b: tb,
				w: tw,
			},
			order: tw.Order,
			want:  tw,
		},
		{
			name:    "Non-existing withdrawal",
			order:   tw.Order,
			wantErr: aerror.WithdrawalFind,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initWithdrawalRepo(t, tt.stored)
			defer r.db.Close()

			eq := mock.ExpectQuery(regexp.QuoteMeta(WithdrawalFind)).WithArgs(tt.order)
			if tt.stored.w.User != "" {
				rows := sqlmock.NewRows([]string{"order", "sum", "processed_at", "user"}).
					AddRow(tt.want.Order, tt.want.Sum, tt.want.ProcessedAt, tt.want.User)
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(aerror.NewError(aerror.WithdrawalFind, sql.ErrNoRows))
			}

			got, err := r.Find(tt.order)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBWithdrawal_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		stored  storedData
		want    []models.Withdrawal
		wantErr string
	}{
		{
			name: "Existing withdrawal",
			stored: storedData{
				b: tb,
				w: tw,
			},
			user: tw.User,
			want: []models.Withdrawal{tw},
		},
		{
			name:    "Non-existing withdrawal",
			user:    tw.User,
			wantErr: aerror.WithdrawalFindAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initWithdrawalRepo(t, tt.stored)
			defer r.db.Close()

			eq := mock.ExpectQuery(regexp.QuoteMeta(WithdrawalFindAll)).WithArgs(tt.user)
			if tt.stored.w.User != "" {
				rows := sqlmock.NewRows([]string{"order", "sum", "processed_at", "user"})
				for _, w := range tt.want {
					rows.AddRow(w.Order, w.Sum, w.ProcessedAt, w.User)
				}
				eq.WillReturnRows(rows)
			} else {
				eq.WillReturnError(aerror.NewError(aerror.WithdrawalFindAll, sql.ErrNoRows))
			}

			got, err := r.FindAll(tt.user)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestNewDBWithdrawalStorage(t *testing.T) {
	db, mock := getMock(t)
	tests := []struct {
		name    string
		db      *sql.DB
		want    DBWithdrawal
		wantErr string
	}{
		{
			name: "Create storage",
			db:   db,
			want: DBWithdrawal{db},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock.ExpectExec(regexp.QuoteMeta(WithdrawalTableCreate)).WillReturnResult(sqlmock.NewResult(1, 1))
			got, err := NewDBWithdrawalStorage(tt.db)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func expectBalanceSet(mock sqlmock.Sqlmock, b models.Balance) {
	mock.ExpectExec(regexp.QuoteMeta(BalanceSet)).
		WithArgs(b.User, b.Current, b.Withdrawn).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func expectWithdrawalAdd(mock sqlmock.Sqlmock, s storedData, duplicate bool) {
	w := s.w
	b := s.b

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(BalanceGet)).WithArgs(w.User).
		WillReturnRows(sqlmock.NewRows([]string{"user", "current", "withdrawn"}).AddRow(b.User, b.Current, b.Withdrawn))

	ee := mock.ExpectExec(regexp.QuoteMeta(WithdrawalAdd)).WithArgs(w.Order, w.Sum, w.ProcessedAt, w.User)
	if duplicate {
		ee.WillReturnError(aerror.NewError(aerror.WithdrawalAdd, sql.ErrNoRows))
		mock.ExpectRollback()
	} else {
		ee.WillReturnResult(sqlmock.NewResult(1, 1))
		expectBalanceSet(mock, models.Balance{
			Current:   b.Current - w.Sum,
			Withdrawn: b.Withdrawn + w.Sum,
			User:      w.User,
		})
		mock.ExpectCommit()
	}
}

func initWithdrawalRepo(t *testing.T, init storedData) (DBWithdrawal, sqlmock.Sqlmock) {
	db, mock := getMock(t)
	r := DBWithdrawal{db}
	b := DBBalance{db}
	if init.w.User != "" {
		expectBalanceSet(mock, init.b)
		if err := b.Set(init.b); err != nil {
			t.Fatal(err)
		}

		expectWithdrawalAdd(mock, init, false)
		if err := r.Add(init.w); err != nil {
			t.Fatal(err)
		}
	}
	return r, mock
}
