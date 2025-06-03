package test

import (
	"BanjirEWS/user"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	repo := user.NewRepository(db.GetDB())

	newUser, err := repo.RegisterUser(context.Background(), &user.User{Username: "John", Email: "john@example.com"})
	assert.Nil(t, err)
	assert.NotNil(t, *newUser)
	assert.Equal(t, "John", newUser.Username)
	assert.Equal(t, "john@example.com", newUser.Email)
}

func TestGetUserByEmail(t *testing.T) {
	repo := user.NewRepository(db.GetDB())

	existingUser, err := repo.GetUserByEmail(context.Background(), "johndoe@example.com")
	assert.Nil(t, err)
	assert.NotNil(t, existingUser)
	assert.Equal(t, "John Doe", existingUser.Username)
	assert.Equal(t, "johndoe@example.com", existingUser.Email)
}

func TestUpdatePassword(t *testing.T) {
	repo := user.NewRepository(db.GetDB())
	err := repo.UpdatePassword(context.Background(), "surya@gmail.com", "newPassword")
	assert.Nil(t, err)
}

func TestUpdateProfile(t *testing.T) {
	repo := user.NewRepository(db.GetDB())
	err := repo.UpdateProfile(context.Background(), &user.User{Username: "John Doe", Email: "johndoe@example.com"}, 117, "account")
	assert.Nil(t, err)
}

func TestCreateUserAdmin(t *testing.T) {
	repo := user.NewRepository(db.GetDB())

	err := repo.CreateUser(context.Background(), &user.User{Username: "Asep", Email: "a@example.com", Password: "rahasia"})
	assert.Nil(t, err)
}

func TestDeleteUser(t *testing.T) {
	repo := user.NewRepository(db.GetDB())

	err := repo.DeleteUser(context.Background(), 116)
	assert.Nil(t, err)
}

func TestGetUsers(t *testing.T) {
	repo := user.NewRepository(db.GetDB())

	users, err := repo.GetUsers(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, users)
}

func TestGetUserById(t *testing.T) {
	repo := user.NewRepository(db.GetDB())

	user, err := repo.GetUserById(context.Background(), 117)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Username)
	assert.Equal(t, "johndoe@example.com", user.Email)
}
