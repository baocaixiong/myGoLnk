package golnk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type Config map[string]interface{}

func (cfg *Config) getConfigRecursive(configName string) interface{} {
	if configName == "" {
		return nil
	}
	keys := strings.Split(configName, ".")
	// beego/config/json.go
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

func (cfg *Config) String(keys string) string {
	value := cfg.getConfigRecursive(keys)
	return fmt.Sprint(value)
}

func (cfg *Config) StringOr(keys string, def string) string {
	value := cfg.getConfigRecursive(keys)
	if value != nil {
		if v, ok := value.(string); ok {
			return v
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (cfg *Config) Int(keys string) (int, error) {
	value := cfg.getConfigRecursive(keys)
	if value != nil {
		if v, ok := value.(int); ok {
			return int(v), nil
		} else {
			return 0, errors.New("not int value")
		}
	} else {
		return 0, errors.New("not exist key: " + keys)
	}
}

func (cfg *Config) IntOr(keys string, def int) int {
	value, err := cfg.Int(keys)
	if err != nil {
		return def
	}
	return value
}

func (cfg *Config) Float(keys string) (float64, error) {
	value := cfg.getConfigRecursive(keys)
	if value != nil {
		if v, ok := value.(float64); ok {
			return v, nil
		} else {
			return 0.0, errors.New("not float64 value")
		}
	} else {
		return 0.0, errors.New("not exist key: " + keys)
	}
}

func (cfg *Config) FloatOr(keys string, def float64) float64 {
	value, err := cfg.Float(keys)
	if err != nil {
		return def
	}
	return value
}

func (cfg *Config) RawValue(keys string) (interface{}, error) {
	value := cfg.getConfigRecursive(keys)
	if value != nil {
		return value, nil
	} else {
		return nil, errors.New("not exist key" + keys)
	}
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
