package conf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	Server ServerConf `yaml:"server"`
}

type ServerConf struct {
	Port string `yaml:"port"`
}

var Settings Conf

func init() {
	temp, err := Readfile("./conf/config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	Settings = temp
}

func Readfile(filePath string) (Conf, error) {
	var settings Conf
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return settings, err
	}
	//yaml文件内容影射到结构体中
	err1 := yaml.Unmarshal(yamlFile, &settings)
	if err1 != nil {
		return settings, err1
	}
	return settings, nil
}
