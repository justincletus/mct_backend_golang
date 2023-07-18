package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var SECRET string

func Config() (map[string]string, error) {

	var data map[string]string
	//var datasource string
	var dataByte []byte

	dataByte, err := ioutil.ReadFile("local.json")
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	if err = json.Unmarshal(dataByte, &data); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	SECRET = data["app_secret"]

	return data, nil

}

func GetAppSecret() string {
	_, err := Config()
	if err != nil {
		return ""
	}

	return SECRET
}
