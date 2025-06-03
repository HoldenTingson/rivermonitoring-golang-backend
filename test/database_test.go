package test

import (
	"BanjirEWS/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db = database.OpenDB()

func TestConnection(t *testing.T) {
	d := database.OpenDB()
	assert.NotNil(t, d)
}
