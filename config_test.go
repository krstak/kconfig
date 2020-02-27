package kconfig

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/krstak/testify"
)

func TestLoad(t *testing.T) {
	_, err := Load("./config-test.yaml")

	testify.Nil(t)(err)
}

func TestLoadUnknownFile(t *testing.T) {
	_, err := Load("./unknown.yaml")

	testify.Equal(t)(errors.New("unable to load config: #open ./unknown.yaml: no such file or directory"), err)
}

func TestLoadWrongFormat(t *testing.T) {
	_, err := Load("./config-test-wrong-format.yaml")

	testify.Equal(t)(errors.New("unable to transfer to json: #yaml: line 3: could not find expected ':'"), err)
}

func TestGet(t *testing.T) {
	c, _ := Load("./config-test.yaml")

	testify.Equal(t)(":8080", c.Get("addr"))
	testify.Equal(t)("postgres", c.Get("database.dialect"))
	testify.Equal(t)("", c.Get("unknown"))
}

func TestGetInt(t *testing.T) {
	c, _ := Load("./config-test.yaml")

	testify.Equal(t)(12, c.GetInt("database.timeout"))
	testify.Equal(t)(0, c.GetInt("unknown"))
}

func TestGetSlice(t *testing.T) {
	c, _ := Load("./config-test.yaml")

	testify.Equal(t)([]string{}, c.GetSlice("list-unknow"))
	testify.Equal(t)([]string{"first", "second", "third"}, c.GetSlice("list"))
}

func TestGet_Env(t *testing.T) {
	envGTM, envPass, unset := setEnvVar()
	defer unset()

	file := "./config-test.yaml"

	c, err := Load(file)

	testify.Nil(t)(err)
	testify.Equal(t)(envGTM, c.Get("gtm"))
	testify.Equal(t)(fmt.Sprintf("postgres://dbuser:%s@localhost:5432/dbname?sslmode=disable", envPass), c.Get("database.url"))
}

func setEnvVar() (string, string, func()) {
	unset := func() {
		os.Unsetenv("GTM_ENV")
		os.Unsetenv("DB_PASS")
	}

	envGTM := "GTM-8788"
	envPass := "super-pass"

	os.Setenv("GTM_ENV", envGTM)
	os.Setenv("DB_PASS", envPass)

	return envGTM, envPass, unset
}
