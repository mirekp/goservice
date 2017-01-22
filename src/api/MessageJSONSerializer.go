package main

import "fmt"

type MessageJSON struct {
	Data MessageData `json:"data"`
}

type MessageData struct {
	ID            int                  `json:"id,omitempty"`
	Type          string               `json:"type,omitempty"`
	Attr          MessageAttributes    `json:"attributes"`
	Relationships MessageRelationships `json:"relationships,omitempty"`
}

type MessageAttributes struct {
	Text string   `json:"text,omitempty"`
}

type MessageRelationships struct {
	Author MessageRelationshipsAuthor   `json:"author,omitempty"`
}

type MessageRelationshipsAuthor struct {
	Links MessageRelationshipsAuthorLinks   `json:"links,omitempty"`
}

type MessageRelationshipsAuthorLinks struct {
	Related string   `json:"related,omitempty"`
}


func serializeMessageToJSON(message Message) MessageJSON {
	return MessageJSON{
		Data: MessageData{
			ID: message.ID,
			Type: "message",
			Attr: MessageAttributes{
				Text: message.Text,
			},
			Relationships: MessageRelationships{
				Author: MessageRelationshipsAuthor{
					Links: MessageRelationshipsAuthorLinks{
						Related:   fmt.Sprintf("http://%s:%d%s/%d", HostName, Port, RouteAccounts, message.Author),
					},
				},
			},
		},
	}
}
