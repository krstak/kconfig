package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
)

type Config struct {
	json string
}

func (c Config) Get(key string) string {
	res := gjson.Get(c.json, key)

	if res.Type == gjson.Null {
		return ""
	}

	return res.String()
}

func (c Config) GetInt(key string) int {
	res := gjson.Get(c.json, key)

	if res.Type == gjson.Null {
		return 0
	}

	if res.Type == gjson.Number {
		return int(res.Int())
	}

	return 0
}

func (c Config) GetSlice(key string) []string {
	res := gjson.Get(c.json, key)

	if res.Type == gjson.Null {
		return []string{}
	}

	arr := res.Array()
	sl := make([]string, 0, len(arr))
	for _, i := range arr {
		sl = append(sl, i.String())
	}
	return sl
}

func Load(path string) (Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to load config: #%v", err)
	}
	json, err := yaml.YAMLToJSON(f)
	if err != nil {
		return Config{}, fmt.Errorf("unable to transfer to json: #%v", err)
	}

	return Config{envVar(string(json))}, nil
}

func envVar(json string) string {
	varReg := regexp.MustCompile(`env\[([^\]]+)\]`)
	envs := varReg.FindAllStringSubmatch(json, -1)
	for _, env := range envs {
		json = strings.Replace(json, env[0], os.Getenv(env[1]), -1)
	}

	return json
}
