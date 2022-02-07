package status

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type (
	Code int

	SimpleStatus struct {
		Code      Code   `json:"status_code"`
		Message   string `json:"status_msg"`
		Timestamp string `json:"timestamp"`
	}
)

const (
	OK Code = iota
	Warning
	Critical
	Unknown

	contentType string = "Content-Type"
	appJson     string = "application/json"
)

func NewStatus() *SimpleStatus {
	return &SimpleStatus{Code: Unknown, Message: "initializing", Timestamp: currentTimestamp()}
}

func (s *SimpleStatus) status(code Code, message string) {
	s.Code = code
	s.Message = message
	s.Timestamp = currentTimestamp()
}

func (s *SimpleStatus) Current() SimpleStatus {
	s.Timestamp = currentTimestamp()
	return *s
}

func currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (s *SimpleStatus) Ok(message string) {
	s.status(OK, message)
}

func (s *SimpleStatus) Warn(message string) {
	s.status(Warning, message)
}

func (s *SimpleStatus) Critical(message string) {
	s.status(Critical, message)
}

func (s *SimpleStatus) Unknown(message string) {
	s.status(Unknown, message)
}

/*
Handler is used for a slowly changing status where we want to automatically update the timestamp to the current request time
*/
func (s *SimpleStatus) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add(contentType, appJson)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(s.Current()); err != nil {
		log.Printf("Couldn't encode/write status: %v", err)
	}
}

/*
BackgroundHandler is used when there will be a background process that updates the status,
and we want to see the timestamp of when the background task ran last
*/
func (s *SimpleStatus) BackgroundHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add(contentType, appJson)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(s); err != nil {
		log.Printf("Couldn't encode/write status: %v", err)
	}
}
