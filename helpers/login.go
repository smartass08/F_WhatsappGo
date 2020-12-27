package helpers

import (
	"F_WhatsappGo/utils"
	qrT "github.com/Baozisoftware/qrcode-terminal-go"
	wa "github.com/Rhymen/go-whatsapp"
	"log"
	"strings"
)

var db utils.DB

func Login(wac *wa.Conn, isretry bool){
	if !isretry{
		DelKey()
	}
	qrChan := make(chan string)
	go func() {
		terminal := qrT.New2(qrT.ConsoleColors.BrightBlack, qrT.ConsoleColors.BrightWhite, qrT.QRCodeRecoveryLevels.Low)
		terminal.Get(<-qrChan).Print()
	}()
	sess, err := wac.Login(qrChan)
	if err != nil {
		log.Printf("error during login: %v\n", err.Error())
		if strings.Contains(err.Error(), "qr code scan timed out"){
			if isretry ==  true{
				log.Fatalln("QR didn't got scan 2 times in a row, Exiting..")
			}
			Login(wac, true)
		}
		return
	}
	PushKey(sess)
	return
}


func GetKey() (bool, wa.Session) {
	db.Access(utils.GetDbUrl())
	confirm, key := db.GetKey()
	if confirm != true{
		return false, key
	} else {
			return true, key
	}
}

func DelKey(){
	db.Access(utils.GetDbUrl())
	db.DelKeys() //Deletes all the keys lmfao
}

func PushKey(access wa.Session){
	db.Access(utils.GetDbUrl())
	db.Addkey(access)
}