package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

const (
	preferredHttpDisplayPort = 100
	fallbackHttpDisplayPort  = 10100
	musicServerPort          = 99
)

// HttpDisplayPort is the port the queue/danmaku HTTP server actually bound.
// Windows can use privileged port 100; Linux non-root falls back to 10100.
// Override with BILILINE_HTTP_PORT.
var HttpDisplayPort = preferredHttpDisplayPort

func httpDisplayURL(path string) string {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("http://127.0.0.1:%d%s", HttpDisplayPort, path)
}

func musicServerURL(path string) string {
	if path == "" || path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("http://127.0.0.1:%d%s", musicServerPort, path)
}

func httpDisplayListenPorts() []int {
	if envPort := strings.TrimSpace(os.Getenv("BILILINE_HTTP_PORT")); envPort != "" {
		port, err := strconv.Atoi(envPort)
		if err == nil && port > 0 && port <= 65535 {
			return []int{port}
		}
		slog.Error("invalid BILILINE_HTTP_PORT, using defaults", "value", envPort, "err", err)
	}
	return []int{preferredHttpDisplayPort, fallbackHttpDisplayPort}
}

//go:embed Resource/web/index.html
var QueHtmlFile []byte

//go:embed Resource/web/default.css
var cssFile []byte

//go:embed Resource/web/DmDisplay.html
var DmDisplayHtml []byte

//go:embed Resource/web/js/NoSleep.min.js
var NoSleepJs []byte

// 使用互斥锁保护共享资源
var (
	queueLock sync.Mutex
	dmLock    sync.Mutex
)

var (
	QueueChatChan = make(chan []byte, 50)
	DmChatChan    = make(chan []byte, 50)
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	QueueConnMap = make(map[*websocket.Conn]bool)
	DmConnMap    = make(map[*websocket.Conn]bool)
)

func StartWebServer() {
	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)(WebServer())

	var ln net.Listener
	for _, port := range httpDisplayListenPorts() {
		HttpDisplayPort = port
		_, _ = http.Get(httpDisplayURL("/EXIT"))
		addr := ":" + strconv.Itoa(port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			slog.Error("HTTP 展示服务监听失败", "port", port, "err", err)
			continue
		}
		ln = listener
		if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok {
			HttpDisplayPort = tcpAddr.Port
		}
		break
	}
	if ln == nil {
		slog.Error("HTTP 展示服务启动失败：无可用端口")
		return
	}
	slog.Info("starting HTTP display server", "addr", ln.Addr().String())
	if err := http.Serve(ln, handler); err != nil {
		slog.Error("HTTP 展示服务退出", "err", err)
	}
}

func WebServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/LineWs", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			slog.Error("Websocket upgrade failed", "err", err)
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
				slog.Error("Failed to close connection", "err", err)
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
					slog.Error("Failed to write message", "err", err)
					delete(QueueConnMap, w)
				}
			}

		}
	})

	mux.HandleFunc("/DmWs", func(writer http.ResponseWriter, request *http.Request) {

		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			slog.Error("Websocket Upgrade Err", "err", err)
			return
		}
		DmConnMap[conn] = true

		err = conn.WriteMessage(websocket.TextMessage, []byte("Connected"))
		if err != nil {
			slog.Error("Websocket Write Err", "err", err)
			delete(DmConnMap, conn)
			return
		}

		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				slog.Error("Failed to close connection", "err", err)
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
			Chat := <-DmChatChan
			ConnMapCopy := DmConnMap
			for w := range ConnMapCopy {
				err = w.WriteMessage(websocket.TextMessage, Chat)
				if err != nil {
					slog.Error("Failed to write message", "err", err)
					delete(DmConnMap, w)
				}
			}
		}
	})

	// 静态资源响应

	mux.HandleFunc("/web", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(QueHtmlFile)
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/dm", func(writer http.ResponseWriter, request *http.Request) {

		//debugger 读取

		DmDisplayHtml, err := os.ReadFile("Resource/web/DmDisplay.html")
		if err != nil {
			panic(err)
		}

		_, err = writer.Write(DmDisplayHtml)
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
			slog.Error("Find font err", "err", err)
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

	mux.HandleFunc("/NoSleep.min.js", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(NoSleepJs)
		if err != nil {
			return
		}
	})

	// 静态同步接口

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

	mux.HandleFunc("/EXIT", func(writer http.ResponseWriter, request *http.Request) {
		CloseRoomConnect(AppClient, GameId, WsClient, CloseHeartbeatChan)
	})

	return mux
}
