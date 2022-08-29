package monitor

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/net/netutil"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
	"tunl-cli/cmd/tui"
	"tunl-cli/ui"
)

const threads = 5

type TunlMonitor struct {
	ws      websocket.Upgrader
	httpSrv *http.Server
	conn    []*websocket.Conn
	mu      sync.Mutex
}

func NewTunlMonitor() *TunlMonitor {
	return &TunlMonitor{
		ws: websocket.Upgrader{},
	}
}

func (m *TunlMonitor) connect(w http.ResponseWriter, r *http.Request) {
	m.ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := m.ws.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer func(conn *websocket.Conn) {
		conn.Close()
	}(conn)

	idx := len(m.conn)
	defer func(m *TunlMonitor) {
		m.conn = append(m.conn[:idx], m.conn[idx+1:]...)
	}(m)

	m.conn = append(m.conn, conn)

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

func (m *TunlMonitor) Start(addr string) error {
	monitorUI, err := fs.Sub(ui.MonitorUI, "monitor/build")
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	listener = netutil.LimitListener(listener, threads)

	r := mux.NewRouter()
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
