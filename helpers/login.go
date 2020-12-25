package helpers

import (
	"F_WhatsappGO/utils"
	"fmt"
	qrT "github.com/Baozisoftware/qrcode-terminal-go"
	wa "github.com/Rhymen/go-whatsapp"
	"os"
)

var db utils.DB

func Login(wac *wa.Conn){
	DelKey()
	qrChan := make(chan string)
	go func() {
		terminal := qrT.New()
		terminal.Get(<-qrChan).Print()
	}()
	sess, err := wac.Login(qrChan)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error during login: %v\n", err.Error())
		return
	}
	PushKey(sess)
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