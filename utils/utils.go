package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

const messageRegex string = `http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`

type ConfigJson struct {
	BOT_TOKEN      string   `json:"bot_token"`
	DB_URL         string   `json:"db_url"`
	DB_Name        string   `json:"db_name"`
	DB_Col         string   `json:"db_collection"`
	Channel_id     int      `json:"channel_id"`
	Links_to_check []string `json:"links_to_check"`
	Filter_mode    string   `json:"filter_mode"`
	Blacklist      []string `json:"blacklist_words"`
	Whitelist      []string `json:"whitelist_words"`
	EmailUser      string   `json:"email_username"`
	EmailPass      string   `json:"email_password"`
	EmailLink      string   `json:"email_imap_link"`
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

func GetfilterMode() string {
	if Config.Filter_mode == "" {
		return "default"
	}
	return Config.Filter_mode
}

func GetChannelId() int {
	return Config.Channel_id
}

func GetLinksToCheck() []string {
	return Config.Links_to_check
}

func GetEmail_Username() string {
	return Config.EmailUser
}

func GetEmail_Password() string {
	return Config.EmailPass
}

func GetEmail_Link() string {
	return Config.EmailLink
}

func LinksValid(message string) (bool, []string) {
	var links []string
	if !strings.Contains(message, "Content-Type: text/plain;") {
		return false, nil
	}
	message = strings.Split(strings.Split(message, "Content-Type: text/plain;")[1],"Content-Type: text/html;")[0]
	match := regexp.MustCompile(messageRegex)
	matches := match.FindAllString(message, -1)
	for _, v := range matches {
		for _, b := range GetLinksToCheck() {
			if strings.Contains(strings.ToLower(v), strings.ToLower(b)) != false {
				links = append(links, v)
				log.Println(v)
			}
		}
	}
	if len(links) > 0 {
		return true, links
	}
	return false, nil
}

func MessageValid(message string) (bool, []string) {
	switch {
	case strings.EqualFold(GetfilterMode(), "blacklist"):
		for _, v := range Config.Blacklist {
			if strings.Contains(message, v) != false {
				return false, nil
			}
		}
		return LinksValid(message)
	case strings.EqualFold(GetfilterMode(), "whitelist"):
		for _, v := range Config.Whitelist {
			if strings.Contains(message, v) != false {
				return LinksValid(message)
			}
		}
		return false, nil
	case strings.EqualFold(GetfilterMode(), "default"):
		return LinksValid(message)
	}

	return false, nil
}

func ReverseInts(input []uint32) []uint32 {
	if len(input) == 0 {
		return input
	}
	return append(ReverseInts(input[1:]), input[0])
}
