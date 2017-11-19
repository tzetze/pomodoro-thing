package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Session represents a running Pomodoro session
type Session struct {
	ResponseURL string
	UserName    string
	TaskName    string
}

// SlackMsg represents a message to Slack
type SlackMsg struct {
	Text string `json:"text"`
}

var sessions map[string]Session

func main() {
	sessions = make(map[string]Session)
	http.HandleFunc("/pomodoro-start", handleSlackStartCommand)
	http.ListenAndServe(":8080", nil)
}

func handleSlackStartCommand(w http.ResponseWriter, r *http.Request) {
	if _, sessionExists := sessions[userIdentifier(r)]; sessionExists {
		respondToSlack(w, "Session is running already.")
		return
	}
	respondToSlack(w, fmt.Sprintf("Starting session: %s", r.FormValue("text")))
	go startPomodoro(r)
}

func startPomodoro(r *http.Request) {
	fmt.Println("Pomodoro start")
	sessions[userIdentifier(r)] = Session{
		ResponseURL: r.FormValue("response_url"),
		UserName:    r.FormValue("user_name"),
		TaskName:    r.FormValue("text"),
	}
	<-time.After(25 * time.Minute)
	go endPomodoro(userIdentifier(r))
}

func endPomodoro(userIdentifier string) {
	fmt.Println("Pomodoro end")
	sessionData := sessions[userIdentifier]
	delete(sessions, userIdentifier)
	sendMessageToSlack("Session finished", sessionData)
}

func userIdentifier(r *http.Request) string {
	return r.FormValue("team_id") + r.FormValue("user_id")
}

func sendMessageToSlack(message string, sessionData Session) {
	msg := buildSlackMessage(fmt.Sprintf("%s: %s", message, sessionData.TaskName))
	_, err := http.Post(sessionData.ResponseURL, "application/json", bytes.NewReader(msg))
	if err != nil {
		fmt.Println(err)
	}
}

func respondToSlack(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	msg := buildSlackMessage(message)
	_, err := w.Write(msg)
	if err != nil {
		fmt.Println("Couldn't send response to slack.")
	}
}

func buildSlackMessage(message string) []byte {
	msg := SlackMsg{Text: message}
	slackMessage, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Couldn't marshal response: %v\n", err)
	}
	return slackMessage
}
