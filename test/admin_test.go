package test

import (
	"context"
	"testing"

	"BanjirEWS/admin"

	"github.com/stretchr/testify/assert"
)

func TestGetAdminByUsername(t *testing.T) {
	repo := admin.NewRepository(db.GetDB())

	adminUser, err := repo.GetAdminByUsername(context.Background(), "adminpekanbaru")
	assert.Nil(t, err)
	assert.NotNil(t, adminUser)
	assert.Equal(t, "adminpekanbaru", adminUser.Username)
	assert.Equal(t, "pekanbaru123", adminUser.Password)
}

func TestGetAllUserCount(t *testing.T) {
	repo := admin.NewRepository(db.GetDB())

	count, err := repo.GetAllUserCount(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 7, count)
}
