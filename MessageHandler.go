package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "strconv"
)

/*
{
    data : {
        "id": 1,
        "type": "message",
        "attributes": {
            "text", "this is a message text"
         },
         "links": {
            "self": "http://example.com/messages/1"
         },
         "relationships": {
            "author": {
              "links": {
                  "related": "http://example.com/accounts/accountname"
              }
            }
         }
    }
}
 */

type Message struct {
    ID            int
    Text          string
    Author        string
}

type MessageJSON struct {
    Data MessageData `json:"data"`
}

type MessageData struct {
    ID             int                  `json:"id,omitempty"`
    Type           string               `json:"type,omitempty"`
    Attr           MessageAttributes    `json:"attributes"`
    Relationships  MessageRelationships `json:"relationships,omitempty"`
}

type MessageAttributes struct {
    Text           string   `json:"text,omitempty"`
}

type MessageRelationships struct {
    Author         MessageRelationshipsAuthor   `json:"author,omitempty"`
}

type MessageRelationshipsAuthor struct {
    Links          MessageRelationshipsAuthorLinks   `json:"links,omitempty"`
}

type MessageRelationshipsAuthorLinks struct {
    Related        string   `json:"related,omitempty"`
}

var messages = []Message{}
var lastMessageID = 0

func msgGetAllHandler (writer http.ResponseWriter, _ *http.Request) {

    if len(messages) == 0 {
        writer.WriteHeader(http.StatusNoContent)
        return
    }


    var messageJSONArray []MessageJSON

    for _, message := range messages {
        messageJSON := serializeMessageToJSON(message)
        messageJSONArray = append(messageJSONArray, messageJSON)
    }

    writer.Header().Set("Content-Type", "application/json")
    json.NewEncoder(writer).Encode(messageJSONArray)
}


func msgPostHandler (writer http.ResponseWriter, request *http.Request) {

    // parse JSON request
    var messageJSON MessageJSON
    err := json.NewDecoder(request.Body).Decode(&messageJSON)
    if err != nil || messageJSON.Data.Type != "message" {
        fmt.Printf("Unable to process your request\n")
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    // client-specified message id's are not (yet) supported
    if messageJSON.Data.ID != 0 {
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    // generate a new messageID, create message model object and store it
    lastMessageID++
    message := Message{
        ID: lastMessageID,
        Text: messageJSON.Data.Attr.Text,
        Author: "an-author",
    }
    messages = append(messages, message)

    // return back the 201 status code and created message
    writer.WriteHeader(http.StatusCreated)
    writer.Header().Set("Content-Type", "application/json")
    createdMsgJSON := serializeMessageToJSON(message)
    json.NewEncoder(writer).Encode(createdMsgJSON)
}

func getMsgIDFromRequest(request *http.Request) (int, error) {
    vars := mux.Vars(request)
    stringID, ok := vars["id"]

    if ok == false  {
        return 0, fmt.Errorf("No message-id provided")
    }

    id, err := strconv.Atoi(stringID)
    if err != nil {
        return 0, fmt.Errorf("Bad msg-id format")
    }

    return id, nil
}

func msgGetHandler (writer http.ResponseWriter, request *http.Request) {
    id, err := getMsgIDFromRequest(request)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    // find the message and return it back
    for _, message := range messages {
        if message.ID == id {
            writer.Header().Set("Content-Type", "application/json")
            writer.WriteHeader(http.StatusOK)
            createdMsgJSON := serializeMessageToJSON(message)
            json.NewEncoder(writer).Encode(createdMsgJSON)
            return
        }
    }

    fmt.Printf("Requested a message %d that doesn't exist\n", id)
    writer.WriteHeader(http.StatusNotFound)
}

func msgEditHandler (writer http.ResponseWriter, request *http.Request) {

    id, err := getMsgIDFromRequest(request)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    // parse JSON request
    var messageJSON MessageJSON
    err = json.NewDecoder(request.Body).Decode(&messageJSON)
    if err != nil {
        fmt.Printf("Unable to process your request\n")
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    // find the message and edit it
    for index, message := range messages {
        if message.ID == id {
            if message.Author == "an-author" {
                // OK - found a message and owner matches. Change the text and return the resource
                fmt.Printf("Modifying message %d\n", id)
                message.Text = messageJSON.Data.Attr.Text
                messages[index] = message
                writer.WriteHeader(http.StatusOK)
                return
            } else {
                fmt.Printf("Unauthorized attempt to edit message %d\n", id)
                // attempted to edit a message that doesn't belong to the currently logged-in user
                writer.WriteHeader(http.StatusUnauthorized)
                return
            }
        }
    }

}

func msgDeleteHandler (writer http.ResponseWriter, request *http.Request) {

    id, err := getMsgIDFromRequest(request)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        return
    }

    for index, message := range messages {
        if message.ID == id {
            if message.Author == "an-author" {
                // OK - found a message and owner matches. Delete it.
                fmt.Printf("Deleting message %d\n", id)
                messages = append(messages[:index], messages[index+1:]...)
                writer.WriteHeader(http.StatusNoContent)
                return
            } else {
                fmt.Printf("Unauthorized attempt to delete message %d\n", id)
                // attempted to delete a message not belonging to the currently logged-in user
                writer.WriteHeader(http.StatusUnauthorized)
                return
            }
        }
    }

    fmt.Printf("Requested deleting message %d that doesn't exist\n", id)
    writer.WriteHeader(http.StatusNotFound)
}

func serializeMessageToJSON(message Message) MessageJSON {
    var messageJSON MessageJSON
    messageJSON.Data.ID = message.ID
    messageJSON.Data.Type = "message"
    messageJSON.Data.Attr.Text = message.Text
    messageJSON.Data.Relationships.Author.Links.Related = fmt.Sprintf("http://%s:%d%s/%s", HostName, Port, AccountsRoute, message.Author)
    return messageJSON
}
