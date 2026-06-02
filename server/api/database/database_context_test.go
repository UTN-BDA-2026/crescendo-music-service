package database_test

import (
	"testing"

	"crescendo-api/database"

	"github.com/stretchr/testify/assert"
)

func TestCanConnectToPostgresDB(t *testing.T) {
	db, err := database.NewConnection()

	assert.NoError(t, err)
	assert.NotNil(t, db)

	defer db.Close()

	assert.NoError(t, db.Ping())
}
