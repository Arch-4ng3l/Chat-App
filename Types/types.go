package types

type Message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Msg      string `json:"msg"`
}

type Friend struct {
	User1    string `json:"user1"`
	User2    string `json:"user2"`
	Accepted bool   `json:"accepted"`
}

type AdminRequest struct {
	Password string `json:"password"`
}
