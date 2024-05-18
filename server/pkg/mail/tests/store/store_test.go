package store_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"passwordless-mail-server/pkg/model"
	"passwordless-mail-server/pkg/util"
	"sync"
	"testing"
	"time"
)

var testDatabase util.TestDatabase
var err error

// postgres needs time to create and drop tables
const waitTime = time.Millisecond * 500

func TestMain(m *testing.M) {
	testDatabase, err = util.NewTestDatabase()
	if err != nil {
		log.Fatal(err)
	}
	err = testDatabase.CreateTestTable()
	fmt.Println("create table error", err)
	time.Sleep(waitTime)
	defer testDatabase.DropTestTable()
	code := m.Run()
	os.Exit(code)
}

func mockMail(amount int) []model.MailEntity {
	var mails []model.MailEntity

	for i := 0; i < amount; i++ {
		mail := model.MailEntity{
			Recipient:   fmt.Sprintf("recipient-%d", i+1),
			Sender:      fmt.Sprintf("sender-%d", i+1),
			MailSubject: fmt.Sprintf("subject-%d", i+1),
			Body:        fmt.Sprintf("body-%d", i+1),
		}
		mails = append(mails, mail)
	}

	return mails
}

func insertMails(mails []model.MailEntity, db *sql.DB) error {
	var wg sync.WaitGroup
	queryString := "INSERT INTO mail (recipient, sender, mail_subject, body) VALUES ($1, $2, $3, $4)"
	for _, mail := range mails {
		wg.Add(1)
		go func(mail model.MailEntity) {
			defer wg.Done()
			_, err := db.Exec(queryString, mail.Recipient, mail.Sender, mail.MailSubject, mail.Body)
			if err != nil {
				return
			}
		}(mail)
	}

	wg.Wait() // wait for all goroutines to finish

	return nil
}

func retrieveMails(db *sql.DB) []model.MailEntity {
	var mails []model.MailEntity
	rows, err := db.Query("SELECT * FROM mail")
	if err != nil {
		return []model.MailEntity{}
	}
	defer rows.Close()

	for rows.Next() {
		var mail model.MailEntity
		err := rows.Scan(&mail.ID, &mail.Recipient, &mail.Sender, &mail.MailSubject, &mail.Body, &mail.SentAt)
		if err != nil {
			return []model.MailEntity{}
		}
		mails = append(mails, mail)
	}

	return mails
}
