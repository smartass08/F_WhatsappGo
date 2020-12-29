package helpers

import (
	"F_WhatsappGo/utils"
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"log"
	"strings"
	"time"
	)

var tg_client = Initialise()

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
	var sender_name string

	if !message.Info.FromMe{
		info_channel, err := h.C.GetGroupMetaData(message.Info.RemoteJid)
		if err != nil {
			log.Printf("Can't get information about the message, Reason %s\n",err.Error())
		}
		info := <-info_channel

		//PM
		if strings.Contains(info, "{\"status\":404}") == true{   //Will not return 404 if its a group
			sender_name = fmt.Sprintf("%v | +%v",h.C.Store.Contacts[message.Info.RemoteJid].Notify,
				strings.Split(h.C.Store.Chats[message.Info.RemoteJid].Jid, "@")[0])

		} else { //Group
			var jid string
			if message.Info.Source.Participant == nil {
				jid = message.Info.RemoteJid
			}else{
				jid = *message.Info.Source.Participant
			}
			sender_name = fmt.Sprintf("%v | %v", h.C.Store.Chats[message.Info.RemoteJid].Name, h.C.Store.Contacts[jid].Notify)
		}
		if utils.MessageValid(message.Text) != false{
			if utils.LinksValid(message.Text) != false{
				log.Println(sender_name)
				TG_Send(tg_client, message.Text, sender_name, false)
			}

		}

	}

}