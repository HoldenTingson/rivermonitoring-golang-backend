package test

import (
	"BanjirEWS/river"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRiver(t *testing.T) {
	req := river.River{
		Id:        "R301",
		Latitude:  "123",
		Longitude: "111",
		Location:  "Pekanbaru",
	}

	err := river.Repository.CreateRiver(river.NewRepository(db.GetDB()), context.Background(), &req)
	assert.Nil(t, err)
}

func TestUpdateRiverDetail(t *testing.T) {
	req := river.River{
		Id:        "R301",
		Latitude:  "12334",
		Longitude: "111",
		Location:  "Pekanbaru",
	}

	err := river.Repository.UpdateRiverDetail(river.NewRepository(db.GetDB()), context.Background(), &req, "R300")
	assert.Nil(t, err)
}

func TestDeleteRiver(t *testing.T) {
	err := river.Repository.DeleteRiver(river.NewRepository(db.GetDB()), context.Background(), "R300")
	assert.Nil(t, err)
}

func TestGetRiver(t *testing.T) {
	rivers, err := river.Repository.GetRiver(river.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, rivers)
}

func TestGetRiverId(t *testing.T) {
	ids, err := river.Repository.GetRiverId(river.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, ids)
}

func TestGetRiverById(t *testing.T) {
	r, err := river.Repository.GetRiverById(river.NewRepository(db.GetDB()), context.Background(), "R301")
	assert.Nil(t, err)
	assert.Equal(t, "R301", r.Id)
}

func TestFindRiver(t *testing.T) {
	r, err := river.Repository.FindRiver(river.NewRepository(db.GetDB()), context.Background(), "Pekanbaru")
	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func TestFilterRiver(t *testing.T) {
	rivers, err := river.Repository.FilterRiver(river.NewRepository(db.GetDB()), context.Background(), "location")
	assert.Nil(t, err)
	assert.NotEmpty(t, rivers)
}

func TestGetAllRiverCount(t *testing.T) {
	count, err := river.Repository.GetAllRiverCount(river.NewRepository(db.GetDB()), context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 6, count)
}

func TestGetRiverByStatus(t *testing.T) {
	rivers, err := river.Repository.GetRiverByStatus(river.NewRepository(db.GetDB()), context.Background(), "bahaya")
	assert.Nil(t, err)
	assert.NotEmpty(t, rivers)
}
