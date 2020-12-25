package main

import (
	"F_WhatsappGO/utils"
	"fmt"
	wa "github.com/Rhymen/go-whatsapp"
	"strings"
	"time"
)


func main()  {
	wac, err := wa.NewConnWithOptions(&wa.Options{
		ShortClientName: "F_whatsappGo",
		LongClientName: "F_whatsappGo",
		Timeout: time.Second*15})
	if err != nil {
		panic(err)
	}
	check, sess := utils.GetKey()
	if check != true{
		fmt.Println("No access token found on db, Need to login")
		utils.Login(wac)
	}
	fmt.Println("Got access token, Trying to login")
	_, err = wac.RestoreWithSession(sess)
	if err != nil{
		if strings.Contains(err.Error(), "admin login responded") == true{
			fmt.Println("Access token Expired, Need re-login")
			utils.Login(wac)
		}
	}
	fmt.Printf("login successful")
	time.Sleep(time.Second*2)
	fmt.Printf("\nIs connected ? :%v \n", wac.Info.Connected)

}