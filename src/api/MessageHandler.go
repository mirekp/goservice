package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
)

func msgGetAllHandler(writer http.ResponseWriter, _ *http.Request) {

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

	if json.NewEncoder(writer).Encode(messageJSONArray) != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func msgPostHandler(writer http.ResponseWriter, request *http.Request) {

	// parse JSON request
	var messageJSON MessageJSON
	err := json.NewDecoder(request.Body).Decode(&messageJSON)
	if err != nil || messageJSON.Data.Type != "message" {
		fmt.Printf("Unable to process your request\n")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// client-specified message id's are not supported
	if messageJSON.Data.ID != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// create message
	message, _ := CreateMessage(0, messageJSON.Data.Attr.Text)

	// and return back the 201 status code and created message
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	createdMsgJSON := serializeMessageToJSON(message)
	json.NewEncoder(writer).Encode(createdMsgJSON)
}

func msgGetHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := getMsgIDFromRequest(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	message := GetMessageByID(id)
	if message == nil {
		fmt.Printf("Requested a message %d that doesn't exist\n", id)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Location", fmt.Sprintf("http://%s:%d%s/%d", HostName, Port, RouteMessage, id))
	createdMsgJSON := serializeMessageToJSON(*message)
	if json.NewEncoder(writer).Encode(createdMsgJSON) != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func msgEditHandler(writer http.ResponseWriter, request *http.Request) {

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
	message := GetMessageByID(id)
	if message == nil {
		fmt.Printf("Requested editing message %d that doesn't exist\n", id)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO - enable once session authentication is implemented
	//if message.Author != loggedInUser {
	//    fmt.Printf("Unauthorized attempt to edit message %d\n", id)
	//    // attempted to edit a message that doesn't belong to the currently logged-in user
	//    writer.WriteHeader(http.StatusUnauthorized)
	//    return
	//}

	// OK - found a message and owner matches. Change the text and return the resource
	fmt.Printf("Modifying message %d\n", id)
	message.Text = messageJSON.Data.Attr.Text
	writer.WriteHeader(http.StatusOK)
}

func msgDeleteHandler(writer http.ResponseWriter, request *http.Request) {

	id, err := getMsgIDFromRequest(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	message := GetMessageByID(id)
	if message == nil {
		fmt.Printf("Requested deleting message %d that doesn't exist\n", id)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO - enable once session authentication is implemented
	//if message.Author != loggedInUser {
	//    fmt.Printf("Unauthorized attempt to delete message %d\n", id)
	//    writer.WriteHeader(http.StatusUnauthorized)
	//    return
	//}

	DeleteMessage(id)
	writer.WriteHeader(http.StatusNoContent)
}


// Helper functions

func getMsgIDFromRequest(request *http.Request) (int, error) {
	vars := mux.Vars(request)
	stringID, ok := vars["id"]

	if ok == false {
		return 0, fmt.Errorf("No message-id provided")
	}

	id, err := strconv.Atoi(stringID)
	if err != nil {
		return 0, fmt.Errorf("Bad msg-id format")
	}

	return id, nil
}
