package main

import (
	_ "embed"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed Resource/web/index.html
var QueHtmlFile []byte

//go:embed Resource/web/default.css
var cssFile []byte

var (
	QueueChatChan = make(chan []byte, 50)
	//DmChatChan    = make(chan []byte, 50)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	QueueConnMap = make(map[*websocket.Conn]bool)
	//DmConnMap    = make(map[*websocket.Conn]bool)
)

func StartWebServer() {
	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)(WebServer())
	err := http.ListenAndServe(":100", handler)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func WebServer() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {

		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}

		QueueConnMap[conn] = true

		err = conn.WriteMessage(websocket.TextMessage, []byte("Connected"))

		if err != nil {
			delete(QueueConnMap, conn)
			return
		}

		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				log.Println("Failed to close connection:", err)
				return
			}
		}(conn)

		go func() {
			for {
				_, Message, err := conn.ReadMessage()
				if err != nil {
					return
				}
				switch string(Message) {
				case "ping":
					err := conn.WriteMessage(websocket.TextMessage, []byte("pong"))
					if err != nil {
						return
					}
				}

			}
		}()

		for {
			Chat := <-QueueChatChan
			ConnMapCopy := QueueConnMap
			for w := range ConnMapCopy {
				err = w.WriteMessage(websocket.TextMessage, Chat)
				if err != nil {
					log.Println("Failed to write message:", err)
					delete(QueueConnMap, w)
					return
				}
			}

		}
	})

	mux.HandleFunc("/web", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(QueHtmlFile)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/font.ttf", func(writer http.ResponseWriter, request *http.Request) {
		err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(info.Name(), ".ttf") {
				file, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				_, err = writer.Write(file)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			log.Println(err)
			return
		}
	})

	mux.HandleFunc("/default.css", func(writer http.ResponseWriter, request *http.Request) {
		var found bool
		dir, err := os.ReadDir("./")
		if err != nil {
			return
		}
		for _, file := range dir {
			if strings.HasSuffix(file.Name(), ".css") {
				found = true
				readFile, err := os.ReadFile(file.Name())
				if err != nil {
					return
				}
				_, err = writer.Write(readFile)
			}
		}

		if !found {
			_, err := writer.Write(cssFile)
			if err != nil {
				return
			}
		}
	})

	mux.HandleFunc("/getAllLine", func(writer http.ResponseWriter, request *http.Request) {
		lineJson, err := json.Marshal(line)
		if err != nil {
			return
		}
		_, err = writer.Write(lineJson)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/getLineLength", func(writer http.ResponseWriter, request *http.Request) {
		LineLength := len(line.GuardLine) + len(line.GiftLine) + len(line.CommonLine)
		_, err := writer.Write([]byte(strconv.Itoa(LineLength)))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/getConfig", func(writer http.ResponseWriter, request *http.Request) {
		ConfigJsonByte, err := json.Marshal(globalConfiguration)
		if err != nil {
			return
		}
		_, err = writer.Write(ConfigJsonByte)
		if err != nil {
			return
		}
	})

	return mux
}
