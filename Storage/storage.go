package storage

import (
	"database/sql"
	"fmt"
	"os"
	types "test/Types"

	loginsystem "github.com/Arch-4ng3l/LoginSystem/LoginSystem"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func New() *SQLite {
	os.Mkdir("Files", os.ModePerm)
	db, err := sql.Open("sqlite3", "Files/database.db")
	if err != nil {
		return nil
	}

	query := `
		CREATE TABLE IF NOT EXISTS users (
			username TEXT UNIQUE,
			email TEXT UNIQUE, 
			password TEXT
		)
	`

	query2 := `
		CREATE TABLE IF NOT EXISTS messages(
			send INT8,
			recv TEXT,
			content TEXT
		)
	`

	query3 := `
		CREATE TABLE IF NOT EXISTS friends(
			user1 TEXT, 
			user2 TEXT,
			accepted BOOL,
			CONSTRAINT unique_user_combination UNIQUE (user1, user2)
		)
	`
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec(query2)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec(query3)
	if err != nil {
		fmt.Println(err)
	}
	return &SQLite{
		db: db,
	}
}

func (s *SQLite) CreateNewUser(req *loginsystem.SignUpRequest) bool {
	query := `
		INSERT INTO users (username, email, password)
		VALUES($1, $2, $3)
	`
	if _, err := s.db.Exec(query, req.Name, req.Email, req.Password); err != nil {
		return false
	}
	return true
}
func (s *SQLite) GetUserInformations(req *loginsystem.LoginRequest) *loginsystem.Account {
	query := `
		SELECT * FROM users WHERE username = $1
	`
	rows, err := s.db.Query(query, req.Name)
	if err != nil {
		return nil
	}
	rows.Next()
	defer rows.Close()
	acc := &loginsystem.Account{}
	if err := rows.Scan(&acc.Name, &acc.Email, &acc.Password); err != nil {
		return nil
	}
	return acc
}

func (s *SQLite) SaveMessage(msg *types.Message) error {
	query := `
		INSERT INTO messages (send, recv, content)
		VALUES($1, $2, $3)
	`
	_, err := s.db.Exec(query, msg.Sender, msg.Receiver, msg.Msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) GetMessages(username string) ([]*types.Message, error) {
	query := `
		SELECT * FROM messages WHERE recv=$1 OR send=$1
	`
	rows, err := s.db.Query(query, username)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	arr := []*types.Message{}
	for rows.Next() {

		msg := &types.Message{}
		err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Msg)
		if err != nil {
			return nil, err

		}
		arr = append(arr, msg)
	}
	return arr, nil

}

func (s *SQLite) GetFriends(username string) ([]*types.Friend, error) {

	query := `
		SELECT * FROM friends WHERE (user1=$1 OR user2=$1) AND accepted=true
	`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	arr := []*types.Friend{}
	for rows.Next() {
		friend := &types.Friend{}
		rows.Scan(&friend.User1, &friend.User2, &friend.Accepted)
		arr = append(arr, friend)
	}

	return arr, nil
}

func (s *SQLite) AddFriend(friend *types.Friend) error {
	query := `
		INSERT INTO friends (user1, user2, accepted)
		VALUES($1, $2, $3)
	`
	_, err := s.db.Exec(query, friend.User1, friend.User2, friend.Accepted)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) AcceptFriend(friend *types.Friend) error {
	query := `
		UPDATE friends
		SET accepted = true
		WHERE user1=$1 AND user2=$2
	`
	_, err := s.db.Exec(query, friend.User1, friend.User2)
	if err != nil {
		return err
	}
	return nil

}

func (s *SQLite) GetFriendRequests(username string) ([]*types.Friend, error) {
	query := `
		SELECT * FROM friends 
		WHERE user2=$1 AND accepted=false;
	`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	arr := []*types.Friend{}
	for rows.Next() {
		friend := &types.Friend{}
		err := rows.Scan(&friend.User1, &friend.User2, &friend.Accepted)
		if err != nil {
			return nil, err
		}
		arr = append(arr, friend)
	}
	return arr, nil
}
