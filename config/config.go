package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var SECRET string

func Config() (string, error) {

	var data map[string]string
	var datasource string

	dataByte, err := ioutil.ReadFile("local.json")
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	if err = json.Unmarshal(dataByte, &data); err != nil {
		return "", fmt.Errorf(err.Error())
	}

	SECRET = data["app_secret"]

	datasource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", data["username"], data["password"], data["host"], data["port"], data["database"])

	return datasource, nil

}
