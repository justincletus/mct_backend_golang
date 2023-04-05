package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var SECRET string

func Config() (string, error) {

	var data map[string]string
	var datasource string
	var dataByte []byte

	host, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("Host name not found")
	}

	if host == "vijin" {
		dataByte, err = ioutil.ReadFile("localDev.json")
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}

	} else {
		dataByte, err = ioutil.ReadFile("local.json")
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}

	}

	if err = json.Unmarshal(dataByte, &data); err != nil {
		return "", fmt.Errorf(err.Error())
	}

	SECRET = data["app_secret"]

	datasource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", data["username"], data["password"], data["host"], data["port"], data["database"])

	return datasource, nil

}
