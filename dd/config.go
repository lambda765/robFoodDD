package dd

import (
	"github.com/spf13/viper"
	"log"
)

const _FILE_PATH string = "./conf.yaml"
const _FILE_TYPE string = "yaml"

type ConfigModel struct {
	Users []UserModel
}

type UserModel struct {
	UserName       string `yaml:"userName""`
	Cookie         string `yaml:"cookie"`
	DdmcUid        string `yaml:"ddmcUid"`
	BarkId         string `yaml:"barkId"`
	AddressNum     int    `yaml:"addressNum"`
	PayMethodNum   int    `yaml:"payMethodNum"`
	SettlementMode int    `yaml:"settlementMode"`
}

var configX = &ConfigModel{}

func InitConfigX() {

	viper.SetConfigFile(_FILE_PATH)
	viper.SetConfigType(_FILE_TYPE)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("read config got an err: ", err)
	}

	viper.Unmarshal(configX)

}
func GetConfigX() *ConfigModel {
	return configX
}
