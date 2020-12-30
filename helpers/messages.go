package helpers

import (
	"F_WhatsappGo/utils"
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"log"
	"strings"
	"sync"
	"time"
)

var email = NewMailService(utils.GetEmail_Username(), utils.GetEmail_Password(), utils.GetEmail_Link())
var tg_client = Initialise()
var wg sync.WaitGroup

type WaHandlers struct {
	C *whatsapp.Conn
}

func (h *WaHandlers) HandleError(err error) {
	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.C.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

func (h *WaHandlers) HandleTextMessage(message whatsapp.TextMessage) {
	var sender_info string

	if !message.Info.FromMe {
		info_channel, err := h.C.GetGroupMetaData(message.Info.RemoteJid)
		if err != nil {
			log.Printf("Can't get information about the message, Reason %s\n", err.Error())
		}
		info := <-info_channel

		//PM
		if strings.Contains(info, "{\"status\":404}") == true { //Will not return 404 if its a group
			sender_info = fmt.Sprintf("%v | +%v", h.C.Store.Contacts[message.Info.RemoteJid].Notify,
				strings.Split(h.C.Store.Chats[message.Info.RemoteJid].Jid, "@")[0])

		} else { //Group
			var jid string
			if message.Info.Source.Participant == nil {
				jid = message.Info.RemoteJid
			} else {
				jid = *message.Info.Source.Participant
			}
			sender_info = fmt.Sprintf("%v | %v", h.C.Store.Chats[message.Info.RemoteJid].Name, h.C.Store.Contacts[jid].Notify)
		}
		check, _ := utils.MessageValid(message.Text)
		if check != false {
			log.Printf("New Invite!, From : %v", sender_info)
			TG_Send(tg_client, message.Text, sender_info, false)
		}
	}
}

func EmailMessages() {
	defer wg.Done()
	log.Println("Checking for new emails")
	raw_mesages, err := email.GetNewMessages()
	if err != nil {
		log.Println("Error getting new messages from server : ", err.Error())
	}
	for _, v := range raw_mesages {
		parsed_message, err := email.ParseMail(v)
		if err != nil {
			log.Println("Error parsing the messages : ", err.Error())
		}
		check, links := utils.MessageValid(parsed_message)
		if check == false {
			err = email.MakeUnread(v.SeqNum)
			if err != nil {
				log.Println("Error making the email as unread : ", err.Error())
			}
		} else {
			from := strings.Split(strings.Split(parsed_message, "From: ")[1], "\n")[0]
			subject := "Empty Subject"
			if strings.Contains(parsed_message, "Subject: ") == true{
				subject = strings.Split(strings.Split(parsed_message, "Subject: ")[1], "\n")[0]
			}
			r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
			sender_info := fmt.Sprintf("%v | %v", r.Replace(from), subject)
			var final_links string
			for _, v := range links {
				final_links += fmt.Sprintf("%v\n", v)
			}
			log.Printf("New Invite!, From : %v", sender_info)
			TG_Send(tg_client, final_links, sender_info, false)
		}
	}
}

func EmailCheckService() {
	for {
		wg.Add(1)
		go EmailMessages()
		wg.Wait()
		time.Sleep(time.Second * 1)

	}
}
