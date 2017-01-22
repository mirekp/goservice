package main

import "fmt"

type Message struct {
	ID            int
	Text          string
	Author        int
}

// Data storage
var messages = []Message{}
var lastMessageID = 0

func GetMessageByID(id int) *Message {

	for _, message := range messages {
		if message.ID == id {
			return &message
		}
	}

	return nil
}

func CreateMessage(author int, text string) (Message, error) {
	message := Message{
		ID: lastMessageID,
		Text: text,
		Author: author,
	}

	messages = append(messages, message)
	fmt.Printf("Posting message %d\n", lastMessageID)

	lastMessageID++
	return message, nil
}

func DeleteMessage(id int) error {
	for idx, message := range messages {
		if message.ID == id {
			fmt.Printf("Deleting message %d\n", id)
			messages = append(messages[:idx], messages[idx + 1:]...)
			return nil
		}
	}
	return fmt.Errorf("Unable to find the message")
}
