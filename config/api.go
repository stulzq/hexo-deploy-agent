package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Get(key string) interface{} {
	for i := len(data) - 1; i >= 0; i-- {
		val, ok := data[i][key]
		if ok && val != nil {
			return val
		}
	}

	return nil
}

func GetString(key string) string {
	iv := Get(key)
	if iv == nil {
		return ""
	}

	if inst, ok := iv.(string); ok {
		return inst
	}

	return ""
}

func GetStruct(key string, dest interface{}) error {

	if key == "" {
		return errors.New("key can not be empty")
	}

	iv := Get(key)

	b, err := yaml.Marshal(iv)
	if b, err = yaml.Marshal(iv); err != nil {
		return errors.Wrap(err, "yaml process: marshal err")
	}

	if err := yaml.Unmarshal(b, dest); err != nil {
		return errors.Wrap(err, "yaml process: unmarshal err")
	}

	return nil
}
