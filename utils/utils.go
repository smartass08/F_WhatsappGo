package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type ConfigJson struct {
	BOT_TOKEN		string 		`json:"bot_token"`
	DB_URL			string		`json:"db_url"`
	DB_Name			string		`json:"db_name"`
	DB_Col			string		`json:"db_collection"`
	Filter_mode		string	 	`json:"filter_mode"`
	Blacklist    	[]string 	`json:"blacklist_links"`
	Whitelist		[]string	`json:"whitelist_links"`

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

func GetFilter_mode() string{
	return Config.Filter_mode
}


func MessageValid(message string) bool{
	switch {
	case strings.EqualFold(GetFilter_mode(), "blacklist"):
		for _, v := range Config.Blacklist{
			if strings.Contains(message, v) != false{
				return false
			}
		}
	case strings.EqualFold(GetFilter_mode(), "whitelist"):
		for _, v := range Config.Whitelist{
			if strings.Contains(message, v) != false{
				return true
			}
		}
	}

	return false
}