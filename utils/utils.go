package utils

import (
	"encoding/json"
	"github.com/DusanKasan/parsemail"
	"io/ioutil"
	"log"
	"mvdan.cc/xurls/v2"
	"strings"
)

//const messageRegex string = `(?:(?:https?|http):\/\/)?[\w/\-?=%.]+\.[\w/\-?=%.]+`

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

func LinksValid(message parsemail.Email, mode string, links []string) (bool, []string) {
	switch {
	case mode == "text":
		links = Getallmatches(message.TextBody)

	case mode == "html":
		links = Getallmatches(message.HTMLBody)
		LinksValid(message, "text", links)
	}

	if len(links) > 0 {
		return true, RemoveDuplicatesStrings(links)
	}
	return false, nil
}

func Cleanlinks(links []string) []string {
	var cleaned []string
	for _, v := range links {
		replacer := strings.NewReplacer(">", "", "<", "", "'", "", `"`, "")
		cleaned = append(cleaned, replacer.Replace(v))
	}
	return cleaned
}
func Getallmatches(message string) []string {
	var links []string
	match := xurls.Strict()
	for _, v := range match.FindAllString(message, -1) {
		for _, b := range GetLinksToCheck() {
			if strings.Contains(strings.ToLower(v), strings.ToLower(b)) != false {
				links = append(links, v)
			}
		}
	}
	return links
}

func MessageValid(message parsemail.Email, mode string, defaultMessage string) (bool, []string) {
	var kek []string
	switch {
	case strings.EqualFold(GetfilterMode(), "blacklist"):
		//Check if its a whatsapp message
		if mode == "default" {
			for _, v := range Config.Blacklist {
				//Compare the whatsapp message
				if strings.Contains(defaultMessage, v) != false {
					return false, nil
				}
			}
			//Check if the slice is empty or not
			if len(Getallmatches(defaultMessage)) > 0 {
				return true, nil
			}
		} else {
			//Do normal operation for emails
			for _, v := range Config.Blacklist {
				if strings.Contains(message.TextBody, v) != false {
					return false, nil
				}
			}
			return LinksValid(message, mode, kek)
		}

	case strings.EqualFold(GetfilterMode(), "whitelist"):
		//Check if its a whatsapp message
		if mode == "default" {
			for _, v := range Config.Whitelist {
				//Compare the whatsapp message
				if strings.Contains(defaultMessage, v) != false {
					//Check if the slice is empty or not after confirming the whitelisted word is inside the text body
					if len(Getallmatches(defaultMessage)) > 0 {
						return true, nil
					}
				}
			}
			return false, nil
		}
		for _, v := range Config.Whitelist {
			if strings.Contains(message.TextBody, v) != false {
				return LinksValid(message, mode, kek)
			}
		}
		return false, nil
	case strings.EqualFold(GetfilterMode(), "default"):
		if mode == "default" {
			if len(Getallmatches(defaultMessage)) > 0 {
				return true, nil
			}
			return false, nil
		} else {
			return LinksValid(message, mode, kek)
		}
	}

	return false, nil
}

func ReverseInts(input []uint32) []uint32 {
	if len(input) == 0 {
		return input
	}
	return append(ReverseInts(input[1:]), input[0])
}

func RemoveDuplicatesStrings(elements []string) []string {
	encountered := map[string]bool{}
	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
