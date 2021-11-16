package main

import (
	"log"
	"time"
)

type Message struct {
	ID               int
	Content, Subject string
	To               int
	From             int
	Created          int64
	Seen             bool
}

func (db *Database) FetchMessages(userID int) ([]Message, error) {

	var messages []Message
	rows, err := db.db.Query("select `id`, `subject`, `from`, `content`, `seen`, `created` from `messages` WHERE `to` = ?", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Subject, &msg.From, &msg.Content, &msg.Seen, &msg.Created)
		if err != nil {
			log.Println(err)
			continue
		}

		messages = append(messages, msg)
	}

	err = rows.Err()
	if err != nil {
		return messages, err
	}

	return messages, nil
}

func (db *Database) SendMessage(to int, from int, subject string, content string) error {
	_, err := db.db.Exec("INSERT INTO `messages`(`to`, `from`, `subject`, `content`, `created`) VALUES (?,?,?,?,?)",
		to, from, subject, content, time.Now().Unix())

	return err
}
