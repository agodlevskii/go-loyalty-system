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

var to = models.Order{
	Number:     `1`,
	Status:     models.StatusNew,
	Accrual:    1,
	UploadedAt: time.Now().Format(time.RFC3339),
	User:       `test`,
}

func TestDBOrder_Add(t *testing.T) {
	tests := []struct {
		name    string
		stored  models.Order
		o       models.Order
		want    models.Order
		wantErr string
	}{
		{
			name:    `Existing order`,
			stored:  to,
			o:       to,
			wantErr: aerror.OrderAdd,
		},
		{
			name: `Non-existing order`,
			o:    to,
			want: to,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initOrderRepo(t, tt.stored)
			defer r.db.Close()

			expectOrderAdd(mock, tt.o, tt.stored.User != ``)
			got, err := r.Add(tt.o)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBOrder_Find(t *testing.T) {
	tests := []struct {
		name    string
		number  string
		stored  models.Order
		want    models.Order
		wantErr string
	}{
		{
			name:   `Existing order`,
			number: to.Number,
			stored: to,
			want:   to,
		},
		{
			name:    `Non-existing order`,
			number:  to.Number,
			wantErr: aerror.OrderFind,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initOrderRepo(t, tt.stored)
			defer r.db.Close()

			eq := mock.ExpectQuery(regexp.QuoteMeta(OrderFind)).WithArgs(tt.number)
			if tt.stored.User != `` {
				eq.WillReturnRows(getOrdersAllRows([]models.Order{tt.want}))
			} else {
				eq.WillReturnError(aerror.NewError(aerror.OrderFind, sql.ErrNoRows))
			}

			got, err := r.Find(tt.number)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBOrder_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		stored  models.Order
		want    []models.Order
		wantErr string
	}{
		{
			name:   `Existing order`,
			user:   to.User,
			stored: to,
			want:   []models.Order{to},
		},
		{
			name:    `Non-xisting order`,
			user:    to.User,
			wantErr: aerror.OrderFindAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initOrderRepo(t, tt.stored)
			defer r.db.Close()

			eq := mock.ExpectQuery(regexp.QuoteMeta(OrderFindAll))
			if tt.stored.User != `` {
				eq.WillReturnRows(getOrdersAllRows(tt.want))
			} else {
				eq.WillReturnError(aerror.NewError(aerror.OrderFindAll, sql.ErrNoRows))
			}

			got, err := r.FindAll(tt.user)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestDBOrder_Update(t *testing.T) {
	tests := []struct {
		name    string
		stored  models.Order
		o       models.Order
		want    models.Order
		wantErr string
	}{
		{
			name:   `Existing order`,
			stored: to,
			o:      to,
			want:   to,
		},
		{
			name:    `Non-existing order`,
			o:       to,
			wantErr: aerror.OrderUpdate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, mock := initOrderRepo(t, tt.stored)
			defer r.db.Close()

			ee := mock.ExpectExec(regexp.QuoteMeta(OrderUpdate)).WithArgs(tt.o.Status, tt.o.Accrual, tt.o.Number)
			if tt.stored.User != `` {
				ee.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				ee.WillReturnError(aerror.NewError(aerror.OrderUpdate, sql.ErrNoRows))
			}

			got, err := r.Update(tt.o)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func TestNewDBOrderStorage(t *testing.T) {
	db, mock := getMock(t)
	tests := []struct {
		name    string
		want    DBOrder
		wantErr string
	}{
		{
			name: `Create storage`,
			want: DBOrder{db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock.ExpectExec(regexp.QuoteMeta(OrderTableCreate)).WillReturnResult(sqlmock.NewResult(1, 1))
			got, err := NewDBOrderStorage(db)
			assert.Equal(t, tt.want, got)
			assertError(t, tt.wantErr, err)
		})
	}
}

func expectOrderAdd(mock sqlmock.Sqlmock, o models.Order, duplicate bool) {
	eq := mock.ExpectQuery(regexp.QuoteMeta(OrderAdd)).WithArgs(o.Number, o.Status, o.Accrual, o.UploadedAt, o.User)
	if duplicate {
		eq.WillReturnError(aerror.NewError(aerror.OrderAdd, sql.ErrNoRows))
	} else {
		eq.WillReturnRows(getOrdersAllRows([]models.Order{o}))
	}
}

func initOrderRepo(t *testing.T, init models.Order) (DBOrder, sqlmock.Sqlmock) {
	db, mock := getMock(t)
	r := DBOrder{db}
	if init.User != `` {
		expectOrderAdd(mock, init, false)
		if _, err := r.Add(init); err != nil {
			t.Fatal(err)
		}
	}
	return r, mock
}

func getOrdersAllRows(orders []models.Order) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{`number`, `status`, `accrual`, `uploaded_at`, `user`})
	for _, o := range orders {
		rows.AddRow(o.Number, o.Status, o.Accrual, o.UploadedAt, o.User)
	}
	return rows
}
