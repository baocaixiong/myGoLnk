package golnk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "reflect"
	"strconv"
	"strings"
)

type Config map[string]interface{}

func (cfg *Config) getConfigRecursive(configName string) interface{} {
	if configName == "" {
		return nil
	}
	keys := strings.Split(configName, ".")
	if len(keys) >= 2 {
		val, ok := (*cfg)[keys[0]]
		if !ok {
			return nil
		}
		for _, key := range keys[1:] {
			if v, ok := val.(map[string]interface{}); !ok {
				return nil
			} else if val, ok = v[key]; !ok {
				return nil
			}
		}
		return val
	} else {
		if v, ok := (*cfg)[configName]; ok {
			return v
		}
	}

	return nil
}

func (cfg *Config) String(keys string, def string) string {
	value := cfg.getConfigRecursive(keys)
	return fmt.Sprint(value)
}
func (cfg *Config) Int(keys string, def int) int {
	value := cfg.getConfigRecursive(keys)
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

func (cfg *Config) Float(keys string, def float64) float64 {
	value := cfg.getConfigRecursive(keys)
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return def
	}
	return f
}

func NewConfig(fileAbsPath string) (*Config, error) {
	config := new(Config)
	bytes, e := ioutil.ReadFile(fileAbsPath)
	if e != nil {
		return config, e
	}
	e = json.Unmarshal(bytes, config)
	return config, e
}
