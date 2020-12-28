package helpers

import (
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

func (m *MailService) ParseMail(raw_mail *imap.Message) (string, error) {
	var mail_body string
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
		mail_body += string(buf)
	}
	return mail_body, nil
}

func (m MailService) MakeUnread(seqnum uint32 ) error {
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

func (m *MailService) GetNewMessages() ([]*imap.Message, error) {
	var mails []*imap.Message
	kek, _ := m.service.Select("INBOX", false)
	cri := imap.NewSearchCriteria()
	cri.WithoutFlags = []string{"\\Seen"}
	_, err := m.service.Search(cri)
	if err != nil {
		return mails, err
	}
	from := uint32(1)
	to := kek.Messages
	if kek.Messages > 1 {
		from = kek.Messages - 29
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
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
	return mail_service
}