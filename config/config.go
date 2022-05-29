package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Gddns map[string]Host `yaml:"gddns"`
}

type Host map[string]Params

type Params struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func GetConfig(confFile string) (Config, error) {
	data, err := ioutil.ReadFile(confFile)
	if err == nil {
		// ymlString := string(data)
		// println(ymlString)

		y := Config{}

		err := yaml.Unmarshal([]byte(data), &y)

		if err != nil {
			log.Fatalf("error: %v", err)
			return Config{}, err
		}

		//fmt.Printf("%+v\n", y)
		return y, nil
	}

	return Config{}, err
}
