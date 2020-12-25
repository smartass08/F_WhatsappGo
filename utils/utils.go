package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigJson struct {
	BOT_TOKEN	string `json:"bot_token"`
	DB_URL		string	`json:"db_url"`
	DB_Name		string	`json:"db_name"`
	DB_Col		string	`json:"db_collection"`
}

const ConfigJsonPath string = "config.json"

var Config *ConfigJson = InitConfig()

func InitConfig() *ConfigJson {
	file, err := ioutil.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}

	var Config ConfigJson
	err = json.Unmarshal([]byte(file), &Config)
	if err != nil {
		log.Fatal(err)
	}
	return &Config
}

func GetBotToken() string {
	return Config.BOT_TOKEN
}

func GetDbUrl() string {
	return Config.DB_URL
}

func GetDbCollection() string {
	return Config.DB_Col
}

func GetDbName() string {
	return Config.DB_Name
}