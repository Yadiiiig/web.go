package library

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Endpoint string
}

func BuildSettings() (Settings, error) {
	var s Settings

	file, err := os.Open("settings.yaml")
	if err != nil {
		return s, err
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(b, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}
