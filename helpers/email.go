package helpers

import (
	"F_WhatsappGo/utils"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)

type MailService struct {
	service    *client.Client
	password   string
	username   string
	serverAddr string
}

func (m *MailService) Connect() error {
	c, err := client.DialTLS(m.serverAddr, nil)
	m.service = c
	return err
}

func (m *MailService) Login() error {
	return m.service.Login(m.username, m.password)
}

func (m *MailService) ParseMail(raw_mail *imap.Message) ([][]byte, error) {
	var mail_body [][]byte
	for _, value := range raw_mail.Body {
		length := value.Len()
		buf := make([]byte, length)
		n, err := value.Read(buf)
		if err != nil {
			return mail_body, err
		}
		if n != length {
			return mail_body, errors.New("Didn't read correct length")
		}
		mail_body = append(mail_body, buf)
	}
	return mail_body, nil
}

func (m MailService) MakeUnread(seqnum uint32) error {
	_, err := m.service.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqnum)
	item := imap.FormatFlagsOp(imap.RemoveFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err = m.service.Store(seqSet, item, flags, nil)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (m MailService) MakeRead(seqnum uint32) error {
	_, err := m.service.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqnum)
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err = m.service.Store(seqSet, item, flags, nil)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (m *MailService) GetNewMessages() ([]*imap.Message, error) {
	var mails []*imap.Message
	_, _ = m.service.Select("INBOX", false)
	cri := imap.NewSearchCriteria()
	cri.WithoutFlags = []string{"\\Seen"}
	uids, err := m.service.Search(cri)
	if err != nil {
		return mails, err
	}
	uids = utils.ReverseInts(uids)
	var limited_uids []uint32
	for i, v := range uids {
		if i > 30 {
			break
		}
		limited_uids = append(limited_uids, v)
	}
	seqset := new(imap.SeqSet)
	seqset.AddNum(limited_uids...)
	items := []imap.FetchItem{imap.FetchRFC822}
	messages := make(chan *imap.Message)
	go func() {
		err = m.service.Fetch(seqset, items, messages)
	}()
	for msg := range messages {
		mails = append(mails, msg)
	}
	return mails, err
}

func NewMailService(username, password, server_addr string) *MailService {
	mail_service := &MailService{username: username, password: password, serverAddr: server_addr}
	err := mail_service.Connect()
	if err != nil {
		log.Println("Error connecting to imap server : ", err.Error())
	}
	err = mail_service.Login()
	if err != nil {
		log.Println("Error logging to imap server : ", err.Error())
	}
	return mail_service
}
