package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"sync"

	ws "github.com/gorilla/websocket"
)

type ConnectionMap struct {
	sync.Mutex
	conns map[string]*ws.Conn
}

func (cm *ConnectionMap) Get(username string) *ws.Conn {
	cm.Lock()
	defer cm.Unlock()
	conn, ok := cm.conns[username]
	if !ok {
		return nil
	}
	return conn
}

func (cm *ConnectionMap) GetAllConnections() []string {
	cm.Lock()
	defer cm.Unlock()
	arr := []string{}
	for key := range cm.conns {
		arr = append(arr, key)
	}
	return arr
}

func (cm *ConnectionMap) Set(username string, conn *ws.Conn) error {
	cm.Lock()
	defer cm.Unlock()
	cm.conns[username] = conn
	return nil
}

func (cm *ConnectionMap) Remove(username string) error {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.conns, username)
	return nil
}

type ApiFunction = func(w http.ResponseWriter, r *http.Request) error

func apiToHttpHandler(f ApiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(err)
		}
	}
}

func getLocalIP() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addr {
		ip, ok := a.(*net.IPNet)
		if ok {
			if !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				return ip.IP.To4().String()
			}
		}
	}
	return ""
}

func createNewSecret() string {
	bytes := make([]byte, 128)
	rand.Read(bytes)
	return base64.RawStdEncoding.EncodeToString(bytes)
}

func getJWTSecret() (string, error) {
	if _, err := os.Stat("Files/secret.txt"); err == nil {
		content, err := os.ReadFile("Files/secret.txt")
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	secret := createNewSecret()
	file, err := os.Create("Files/secret.txt")
	if err != nil {
		return "", err
	}
	if _, err := file.Write([]byte(secret)); err != nil {
		return "", err
	}

	return secret, nil
}
