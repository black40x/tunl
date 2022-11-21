package monitor

import (
	"context"
	"encoding/json"
	"github.com/black40x/tunl-cli/cmd/tui"
	"github.com/black40x/tunl-cli/ui"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/net/netutil"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const threads = 5

type TunlMonitor struct {
	ws        websocket.Upgrader
	httpSrv   *http.Server
	conn      map[int]*websocket.Conn
	connCount int
	template  *template.Template
	mu        sync.Mutex
	appConf   appConfig
}

type appConfig struct {
	Port string
	Host string
}

func NewTunlMonitor() *TunlMonitor {
	return &TunlMonitor{
		ws:   websocket.Upgrader{},
		conn: make(map[int]*websocket.Conn),
	}
}

func (m *TunlMonitor) loadTemplate() {
	index, _ := fs.Sub(ui.MonitorUI, "monitor/build")
	m.template, _ = template.ParseFS(index, "index.html")
}

func (m *TunlMonitor) index(w http.ResponseWriter, r *http.Request) {
	m.template.Execute(w, m.appConf)
}

func (m *TunlMonitor) connect(w http.ResponseWriter, r *http.Request) {
	m.ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := m.ws.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	m.mu.Lock()
	m.connCount++
	m.mu.Unlock()

	idx := m.connCount + 1

	defer func(conn *websocket.Conn) {
		conn.Close()
		delete(m.conn, idx)
	}(conn)

	defer func(m *TunlMonitor) {
		delete(m.conn, idx)
	}(m)

	m.conn[idx] = conn

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		} else {
			log.Println("read:", string(message))
		}
	}
}

func (m *TunlMonitor) SendJsonMessage(d interface{}) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.conn {
		if m.conn[i] != nil {
			m.conn[i].WriteMessage(websocket.TextMessage, data)
		}
	}

	return nil
}

func (m *TunlMonitor) Shutdown(ctx context.Context) {
	for i := range m.conn {
		if m.conn[i] != nil {
			m.conn[i].Close()
		}
	}

	if m.httpSrv != nil {
		m.httpSrv.Shutdown(ctx)
	}
}

func (m *TunlMonitor) Start(addr, host, port string) error {
	monitorUI, err := fs.Sub(ui.MonitorUI, "monitor/build")
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	m.loadTemplate()
	m.appConf = appConfig{
		Host: host,
		Port: port,
	}

	listener = netutil.LimitListener(listener, threads)

	r := mux.NewRouter()
	r.HandleFunc("/", m.index)
	r.HandleFunc("/connect", m.connect)
	r.PathPrefix("/").Handler(http.FileServer(http.FS(monitorUI)))

	m.httpSrv = &http.Server{
		Addr:        addr,
		Handler:     r,
		ReadTimeout: time.Second * 30,
		IdleTimeout: time.Second * 30,
	}

	go func() {
		if err := m.httpSrv.Serve(listener); err != nil {
			if err != http.ErrServerClosed {
				tui.PrintError(err)
				os.Exit(1)
			}
		}
	}()

	return nil
}
