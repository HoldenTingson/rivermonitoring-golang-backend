package test

import (
	"BanjirEWS/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDate(t *testing.T) {
	indonesianDate, err := util.FormatIndonesianDate("2006-01-02")
	assert.Nil(t, err)
	assert.Equal(t, "2 Januari 2006", indonesianDate)
}

func TestCheckPassword(t *testing.T) {
	hashedPassword, err := util.HashPassword("holden")
	assert.Nil(t, err)
	err = util.CheckPassword("holden", hashedPassword)
	assert.Nil(t, err)
}
