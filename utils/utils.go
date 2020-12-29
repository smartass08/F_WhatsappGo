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
	Channel_id		int 		`json:"channel_id"`
	Links_to_check	[]string	`json:"links_to_check"`
	Filter_mode		string	 	`json:"filter_mode"`
	Blacklist    	[]string 	`json:"blacklist_words"`
	Whitelist		[]string	`json:"whitelist_words"`

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

func GetfilterMode() string{
	if Config.Filter_mode == ""{
		return "default"
	}
	return Config.Filter_mode
}

func GetChannelId() int{
	return Config.Channel_id
}

func GetLinksToCheck() []string{
	return Config.Links_to_check
}


func LinksValid(msg string) bool{
	for _,v := range GetLinksToCheck(){
		if strings.Contains(msg, v) != false{
			return true
		}
	}
	return false
}

func MessageValid(message string) bool{
	switch {
	case strings.EqualFold(GetfilterMode(), "blacklist"):
		for _, v := range Config.Blacklist{
			if strings.Contains(message, v) != false{
				return false
			}
		}
		return true
	case strings.EqualFold(GetfilterMode(), "whitelist"):
		for _, v := range Config.Whitelist{
			if strings.Contains(message, v) != false{
				return true
			}
		}
		return false
	case strings.EqualFold(GetfilterMode(), "default"):
		return true
	}

	return false
}