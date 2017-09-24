//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//

// [ Geduld ist eine gute Eigenschaft. Aber nicht, wenn es um die Beseitigung ]
// [ von Missständen geht.                                                    ]
// [                                                      -Margaret Thatcher- ]

package app

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const (
	FILENAME = "app.yml"
)

type (
	ConfigService struct {
		LastErr error
	}
)

var (
	BaseDir  string
	Hub      string
	DSN      string
	Gateways []string
)

//
//
//
func NewConfigService() *ConfigService {
	return &ConfigService{}
}

//
//
//
func (d *ConfigService) RetrieveAll() *ConfigService {

	type Config struct {
		Hub      string   `yaml:"hub"`
		DSN      string   `yaml:"dsn"`
		Gateways []string `yaml:"gateways"`
	}

	if d.LastErr != nil {
		return d
	}

	filename := path.Join(BaseDir, FILENAME)

	config := &Config{}
	d.LastErr = readYML(config, filename)

	Hub = config.Hub
	DSN = config.DSN
	Gateways = config.Gateways

	return d
}

//
//
//
func readYML(v interface{}, filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(raw, v)
}
