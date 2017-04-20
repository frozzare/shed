package app

import (
	"testing"

	assert "github.com/frozzare/go-assert"
	"github.com/frozzare/shed/config"
	"github.com/frozzare/shed/repository"
)

func TestApp(t *testing.T) {
	app, err := NewApp(nil)

	assert.Nil(t, app)
	assert.Equal(t, err, ErrInvalidOptions)
}

func TestDomain(t *testing.T) {
	app, _ := NewApp(&Options{
		Config: config.Config{
			Host: "shed.io",
		},
		Repository: repository.Repository{
			Slug: "master",
		},
	})

	assert.Equal(t, app.Host(), "master.shed.io")
}
