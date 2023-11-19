//////////////////////////////////////////////////////////////////////////////////////////////////////////
///// Author   : Latharani M K
///// Purpose  : The requirements for the log ingestor and the query interface to run in default port 3000
///// Date     : 18th-Nov-2023
///// Language : Golang
//////////////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Log represents the log entry format
type Log struct {
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	ResourceID  string    `json:"resourceId"`
	Timestamp   time.Time `json:"timestamp"`
	TraceID     string    `json:"traceId"`
	SpanID      string    `json:"spanId"`
	Commit      string    `json:"commit"`
	Metadata    Metadata  `json:"metadata"`
}

// Metadata represents the metadata field in the log entry
type Metadata struct {
	ParentResourceID string `json:"parentResourceId"`
}

// LogStorage stores logs and provides query functionality
type LogStorage struct {
	logs []Log
	mu   sync.RWMutex
}

// NewLogStorage creates a new LogStorage instance
func NewLogStorage() *LogStorage {
	return &LogStorage{}
}

// Ingest logs a new log entry
func (ls *LogStorage) Ingest(log Log) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.logs = append(ls.logs, log)
}

// Query searches for logs based on provided filters
func (ls *LogStorage) Query(filters map[string]string) []Log {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	var result []Log

	for _, log := range ls.logs {
		if matchesFilters(log, filters) {
			result = append(result, log)
		}
	}

	return result
}

// matchesFilters checks if a log entry matches the provided filters
func matchesFilters(log Log, filters map[string]string) bool {
	for key, value := range filters {
		switch key {
		case "level":
			if log.Level != value {
				return false
			}
		case "message":
			if !strings.Contains(log.Message, value) {
				return false
			}
		case "resourceId":
			if log.ResourceID != value {
				return false
			}
		case "timestamp":
			timestamp, err := time.Parse(time.RFC3339, value)
			if err != nil || log.Timestamp.Before(timestamp) || log.Timestamp.After(timestamp.Add(24*time.Hour)) {
				return false
			}
		case "traceId":
			if log.TraceID != value {
				return false
			}
		case "spanId":
			if log.SpanID != value {
				return false
			}
		case "commit":
			if log.Commit != value {
				return false
			}
		case "metadata.parentResourceId":
			if log.Metadata.ParentResourceID != value {
				return false
			}
		}
	}

	return true
}

func main() {
	logStorage := NewLogStorage()

	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Ingest called")
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var log Log
		err = json.Unmarshal(body, &log)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		logStorage.Ingest(log)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Query called")
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var filters map[string]string
		err = json.Unmarshal(body, &filters)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		logs := logStorage.Query(filters)

		response, err := json.Marshal(logs)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	fmt.Println("Hi Dyte , Log Ingestor is running on Port :3000...")
	http.ListenAndServe(":3000", nil)
}


