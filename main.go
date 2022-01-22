package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	traceID := strings.Split(r.Header.Get("X-Cloud-Trace-Context"), "/")[0]

	// application log
	logInfo("This is 1st application logs", traceID)
	logInfo("This is 2nd application logs", traceID)
	logInfo("This is 3rd application logs", traceID)
	logInfo("This is 4th application logs", traceID)
	logInfo("This is 5th application logs", traceID)

	// application error log
	logError("This is 1st application error logs", traceID)
	logError("This is 2nd application error logs", traceID)
	logError("This is 3rd application error logs", traceID)

	// response
	fmt.Fprint(w, "Hello, World!")
}

type ApplicationLog struct {
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
	Trace    string    `json:"logging.googleapis.com/trace"`
}

// 本来はctxを渡して、ctxからTraceIDを取るようにしたほうがいいんだと思う。知らんけど。
func logInfo(msg, traceID string) {
	appLog := NewApplicationLog("INFO", msg, traceID)
	json.NewEncoder(os.Stdout).Encode(appLog)
}

func logError(msg, traceID string) {
	appLog := NewApplicationLog("ERROR", msg, traceID)
	json.NewEncoder(os.Stderr).Encode(appLog)
}

func NewApplicationLog(severity, msg, traceID string) *ApplicationLog {
	return &ApplicationLog{
		Severity: severity,
		Message:  msg,
		Time:     time.Now(),
		Trace:    fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), traceID),
	}
}
