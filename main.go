package main

import (
	"F_WhatsappGo/helpers"
	wa "github.com/Rhymen/go-whatsapp"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func issalive(wac *wa.Conn) {
	defer wg.Done()
	for {
		pong, err := wac.AdminTest()
		log.Printf("Is it alive : %v\n", pong)
		if !pong || err != nil {
			if strings.Contains(err.Error(), "connection timed out") {
				log.Printf("Client seems to be disconnected, Please check your phone")
			}
		}
		time.Sleep(time.Second * 5)

	}
}

func main() {
	wg.Add(1)
	go helpers.EmailCheckService()
	//Check latest webui client verion
	v, err := wa.CheckCurrentServerVersion()
	if err != nil {
		panic(err)
	}
	log.Printf("\nLatest client version is := %v.%v.%v\n", v[0], v[1], v[2])
	wac, err := wa.NewConnWithOptions(&wa.Options{
		ShortClientName: "F_whatsappGo",
		LongClientName:  "F_whatsappGo",
		Timeout:         time.Second * 15})
	if err != nil {
		panic(err)
	}
	wac.SetClientVersion(v[0], v[1], v[2])

	check, sess := helpers.GetKey()
	if check != true {
		log.Println("No access token found on db, Need to login")
		helpers.Login(wac, false)
	}
	log.Println("Got access token, Trying to login")
	_, err = wac.RestoreWithSession(sess)
	if err != nil {
		if strings.Contains(err.Error(), "admin login responded") == true {
			log.Println("Access token Expired, Need re-login")
			helpers.Login(wac, false)
		}
	}
	log.Println("login successful")
	time.Sleep(time.Second * 2)
	wac.AddHandler(&helpers.WaHandlers{C: wac})
	wg.Add(1)
	go issalive(wac)

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
