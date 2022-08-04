package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-loyalty-system/internal/aerror"
	"reflect"
	"testing"
)

func TestNewDBRepo(t *testing.T) {
	type args struct {
		url    string
		driver string
	}
	tests := []struct {
		name    string
		args    args
		want    Repo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDBRepo(tt.args.url, tt.args.driver)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDBRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBRepo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertError(t *testing.T, label string, err *aerror.AppError) {
	assert.Equal(t, label == ``, err == nil)
	if err != nil {
		assert.Equal(t, label, err.Label)
	}
}

func getMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	return db, mock
}
