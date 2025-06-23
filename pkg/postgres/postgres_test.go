package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const validDSN = "postgres://postgres:password@localhost:5432/metrics"
const invalidDSN = "postgres://invalid:invalid@localhost:5432/invalid"

func TestPostgres_New_Success(t *testing.T) {
	db, err := New(validDSN)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.Pool)
	assert.NotNil(t, db.Builder)

	// Проверка подключения
	err = db.Pool.Ping(context.Background())
	assert.NoError(t, err)

	db.Close()
}

func TestPostgres_New_Failure(t *testing.T) {
	invalidDSN := "postgres://invalid_user:wrong_password@localhost:5432/wrong_db"

	db, err := New(invalidDSN, ConnAttempts(1))
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Проверка, что ping возвращает ошибку
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	pingErr := db.Pool.Ping(ctx)
	assert.Error(t, pingErr)
}

func TestPostgres_Close(t *testing.T) {
	db, err := New(validDSN)
	assert.NoError(t, err)

	db.Close()

	// Проверка повторного закрытия (не должно паниковать)
	assert.NotPanics(t, func() {
		db.Close()
	})
}

func TestPostgresOptions(t *testing.T) {
	pg := &Postgres{}
	MaxPoolSize(10)(pg)
	ConnAttempts(5)(pg)

	assert.Equal(t, 10, pg.maxPoolSize)
	assert.Equal(t, 5, pg.connAttempts)
}
