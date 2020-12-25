package main

import (
	"F_WhatsappGO/helpers"
	"F_WhatsappGO/utils"
	wa "github.com/Rhymen/go-whatsapp"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)


func main() {
	//Check latest webui client verion
	v, err := wa.CheckCurrentServerVersion()
	if err != nil {
		panic(err)
	}
	log.Printf("\nLatest client version is := %v.%v.%v\n", v[0],v[1],v[2])
	wac, err := wa.NewConnWithOptions(&wa.Options{
		ShortClientName: "F_whatsappGo",
		LongClientName: "F_whatsappGo",
		Timeout: time.Second*15})
	if err != nil {
		panic(err)
	}
	wac.SetClientVersion(v[0],v[1],v[2])

	check, sess := utils.GetKey()
	if check != true{
		log.Println("No access token found on db, Need to login")
		helpers.Login(wac)
	}
	log.Println("Got access token, Trying to login")
	_, err = wac.RestoreWithSession(sess)
	if err != nil{
		if strings.Contains(err.Error(), "admin login responded") == true{
			log.Println("Access token Expired, Need re-login")
			helpers.Login(wac)
		}
	}
	log.Println("login successful")
	time.Sleep(time.Second*2)
	wac.AddHandler(&helpers.WaHandlers{wac})

	//test
	pong, err := wac.AdminTest()
	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	//Disconnect safe
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down now.")
	_, err = wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}
}