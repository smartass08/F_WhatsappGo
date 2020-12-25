package helpers

import (
	"F_WhatsappGO/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func Initialise() *gotgbot.Updater {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zap.InfoLevel))
	l := logger.Sugar()
	updater, err := gotgbot.NewUpdater(logger, utils.GetBotToken())
	l.Info("Got Updater")
	updater.UpdateGetter = ext.BaseRequester{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 65,
		},
		ApiUrl: ext.ApiUrl,
	}
	updater.Bot.Requester = ext.BaseRequester{Client: http.Client{Timeout: time.Second * 65}}
	if err != nil {
		l.Fatalw("failed to start updater", zap.Error(err))
	}
	return updater
}

func TG_Send(client *gotgbot.Updater ,content string, contactInfo string, retry bool) {
	if retry {
		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(60-30) + 30
		log.Println("Sleeping for : ", random)
		time.Sleep(time.Duration(random) * time.Second)
	}
	message := fmt.Sprintf("<b>New Invite Arrived!</b>\n%v\n\n%v", contactInfo, content)
	send := client.Bot.NewSendableMessage(utils.GetChannelId(), message)
	send.ParseMode = "HTML"
	send.DisableWebPreview = true
	_, err := send.Send()
	if err != nil {
		log.Printf("%v\n", err)
		if strings.Contains(err.Error(), "Too Many Requests") {
			TG_Send(client, content, contactInfo, true)
		}
	}
}
