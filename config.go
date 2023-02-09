package kconfig

import (
	"encoding/json"
	"errors"
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

func (c Config) GetBool(key string) bool {
	res := gjson.Get(c.json, key)

	if res.Type == gjson.Null {
		return false
	}

	if res.Type == gjson.True {
		return res.Bool()
	}

	return false
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

func (c Config) GetX(key string, v interface{}) error {
	res := gjson.Get(c.json, key)

	if res.Type == gjson.Null {
		return errors.New("node not found")
	}

	return json.Unmarshal([]byte(res.Raw), v)
}

func LoadContent(content []byte) (Config, error) {
	json, err := yaml.YAMLToJSON(content)
	if err != nil {
		return Config{}, fmt.Errorf("unable to transfer to json: #%v", err)
	}

	return Config{envVar(string(json))}, nil
}

func Load(path string) (Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to load config: #%v", err)
	}
	return LoadContent(f)
}

func envVar(json string) string {
	varReg := regexp.MustCompile(`env\[([^\]]+)\]`)
	envs := varReg.FindAllStringSubmatch(json, -1)
	for _, env := range envs {
		key := env[1]
		def := ""
		if strings.Contains(key, "|") {
			parts := strings.Split(key, "|")
			key = parts[0]
			def = parts[1]
		}

		value, present := os.LookupEnv(key)
		if !present {
			value = def
		}
		json = strings.Replace(json, env[0], value, -1)
	}

	return json
}
