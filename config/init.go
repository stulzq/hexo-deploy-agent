package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	data []map[string]interface{}
)

func init() {
	configFile := "conf/config.yml"

	err := readFileData(configFile, &data)
	if err != nil {
		log.Fatalf("config init failed %v", err)
	}
}

func readFileData(fileName string, data *[]map[string]interface{}) error {
	confFileStream, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	fileData := make(map[string]interface{})
	err = yaml.Unmarshal(confFileStream, &fileData)

	if err != nil {
		return errors.Wrap(err, "Unmarshal err")
	}

	//展开配置
	item := make(map[string]interface{})
	*data = append(*data, item)
	for ck, cv := range fileData {
		item[ck] = cv
		switch inst := cv.(type) {
		case map[string]interface{}:
			expandFileConfig(item, ck, inst)
		}
	}

	return nil
}

func expandFileConfig(data map[string]interface{}, prefix string, current map[string]interface{}) {
	for ck, cv := range current {
		key := fmt.Sprintf("%s:%s", prefix, ck)
		data[key] = cv
		switch inst := cv.(type) {
		case map[string]interface{}:
			expandFileConfig(data, key, inst)
		}
	}
}
