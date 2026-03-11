package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	logFile *os.File
	mu      sync.Mutex
)

func main() {
	var err error
	// Open or create the master execution log
	logFile, err = os.OpenFile("../master_execution.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	http.HandleFunc("/log", handleLog)
	fmt.Println(">> SomaOS Master Logger: WebSocket listening on :8082/log")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		writeLog(fmt.Sprintf("[BROWSER] %s", string(message)))
	}
}

func writeLog(msg string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(msg)
	logFile.WriteString(msg + "\n")
}
