package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	storage "test/Storage"
	types "test/Types"

	loginsystem "github.com/Arch-4ng3l/LoginSystem/LoginSystem"
	color "github.com/TwiN/go-color"
	ws "github.com/gorilla/websocket"
)

type APIServer struct {
	login    *loginsystem.LoginSystem
	store    *storage.SQLite
	conns    *ConnectionMap
	shutdown chan bool
}

var upgrader = ws.Upgrader{}

func New(store *storage.SQLite) *APIServer {
	return &APIServer{
		store: store,
		conns: &ConnectionMap{
			conns: make(map[string]*ws.Conn),
		},
		shutdown: make(chan bool),
	}
}

func (s *APIServer) Run() {

	errChan := make(chan loginsystem.ErrorStruct, 20)
	ls := loginsystem.NewLoginSystem(":3000", "/api", s.store, errChan)

	secret, err := getJWTSecret()
	if err != nil {
		log.Println(err)
		return
	}
	s.login = ls

	s.login.Run(secret)

	ip := getLocalIP()

	http.HandleFunc("/api/fr", apiToHttpHandler(s.handleAddFriend))
	http.HandleFunc("/api/fr/accept", apiToHttpHandler(s.handleAcceptFriend))
	http.HandleFunc("/api/ws", apiToHttpHandler(s.handleWebSocket))
	http.HandleFunc("/api/admin/login", apiToHttpHandler(s.handleAdmin))
	http.HandleFunc("/api/admin/shutdown", apiToHttpHandler(s.handleShutdown))

	fmt.Println("Server Läuft Auf > " + color.Bold + color.Green + ip + ":3000" + color.Reset)
	fmt.Println("Fenster nicht Schließen")

	dir := http.Dir("Website")
	fileServer := http.FileServer(dir)
	http.Handle("/", fileServer)
	go http.ListenAndServe(":3000", nil)
	<-s.shutdown
	log.Println("Shutdown")
}

func (s *APIServer) handleShutdown(w http.ResponseWriter, r *http.Request) error {
	s.shutdown <- true
	return nil
}

func (s *APIServer) handleAdmin(w http.ResponseWriter, r *http.Request) error {
	req := &types.AdminRequest{}
	json.NewDecoder(r.Body).Decode(req)

	if req.Password != "admin" {
		return nil
	}

	conns := s.conns.GetAllConnections()

	return json.NewEncoder(w).Encode(map[string][]string{"conns": conns})
}

func (s *APIServer) handleAddFriend(w http.ResponseWriter, r *http.Request) error {
	req := &types.Friend{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	req.Accepted = false
	if s.store.GetUserInformations(&loginsystem.LoginRequest{Name: req.User2}) == nil {
		return fmt.Errorf("The Person U've tried to add doesnt exist")
	}
	err := s.store.AddFriend(req)
	if err != nil {
		return err
	}
	if conn := s.conns.Get(req.User2); conn != nil {
		conn.WriteJSON(req)
	}

	return nil
}

func (s *APIServer) handleAcceptFriend(w http.ResponseWriter, r *http.Request) error {
	req := &types.Friend{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if err := s.store.AcceptFriend(req); err != nil {
		return err
	}
	if conn := s.conns.Get(req.User1); conn != nil {
		req.Accepted = true
		conn.WriteJSON(req)
	}

	return nil
}

func (s *APIServer) handleWebSocket(w http.ResponseWriter, r *http.Request) error {
	acc := s.login.AuthWithJWT(r)
	if acc == nil {
		conn, _ := upgrader.Upgrade(w, r, nil)
		conn.Close()
		return fmt.Errorf("Not A Valid Token")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := s.conns.Set(acc.Name, conn); err != nil {
		return err
	}
	errCh := make(chan error)

	go s.handleFetchMessage(conn, acc)

	go s.handleSendMessage(conn, acc, errCh)

	err = <-errCh
	s.conns.Remove(acc.Name)
	conn.Close()
	_, ok := err.(*ws.CloseError)
	if ok {
		return nil
	}
	return err
}

func (s *APIServer) handleSendMessage(conn *ws.Conn, acc *loginsystem.Account, errCh chan error) {
	msg := &types.Message{}

	for {
		if err := conn.ReadJSON(msg); err != nil {
			errCh <- err
		}
		msg.Sender = acc.Name
		if err := s.store.SaveMessage(msg); err != nil {
			errCh <- err
		}
		if c := s.conns.Get(msg.Receiver); c != nil {
			if err := c.WriteJSON(msg); err != nil {
				errCh <- err
			}
		}
	}
}
func (s *APIServer) handleFetchMessage(conn *ws.Conn, acc *loginsystem.Account) {
	friends, err := s.store.GetFriends(acc.Name)
	if err != nil || friends == nil {
		return
	}
	for _, friend := range friends {
		conn.WriteJSON(friend)
	}
	friendReqs, err := s.store.GetFriendRequests(acc.Name)
	if err != nil {
		log.Println(err)
		return
	}
	for _, req := range friendReqs {
		conn.WriteJSON(req)
	}

	msgs, err := s.store.GetMessages(acc.Name)
	if err != nil {
		return
	}
	for _, msg := range msgs {
		conn.WriteJSON(msg)
	}
}
