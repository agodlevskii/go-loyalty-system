package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go-loyalty-system/internal/aerror"
	"testing"
)

func assertError(t *testing.T, label string, err *aerror.AppError) {
	assert.Equal(t, label == "", err == nil)
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
