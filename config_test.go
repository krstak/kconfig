package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/krstak/testify"
)

func TestLoad(t *testing.T) {
	defer func() {
		os.Unsetenv("GTM_ENV")
		os.Unsetenv("DB_PASS")
	}()

	envGTM := "GTM-8788"
	envPass := "super-pass"

	os.Setenv("GTM_ENV", envGTM)
	os.Setenv("DB_PASS", envPass)

	file := "./config-test.yaml"

	c, err := Load(file)

	testify.Nil(t)(err)
	testify.Equal(t)(":8080", c.Get("addr"))
	testify.Equal(t)(envGTM, c.Get("gtm"))
	testify.Equal(t)("postgres", c.Get("database.dialect"))
	testify.Equal(t)(fmt.Sprintf("postgres://dbuser:%s@localhost:5432/dbname?sslmode=disable", envPass), c.Get("database.url"))
	testify.Equal(t)(12, c.GetInt("database.timeout"))
	testify.Equal(t)("", c.Get("unknown"))
	testify.Equal(t)(0, c.GetInt("unknown"))
	testify.Equal(t)([]string{}, c.GetSlice("list-unknow"))
	testify.Equal(t)([]string{"first", "second", "third"}, c.GetSlice("list"))
}
