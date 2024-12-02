package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type DatabaseConnection struct {
	initialized bool
	db          *sql.DB
}

var _DBFILE = "./db_files/chat_query.db"

func GetDbFileLocation() string {
	return _DBFILE
}
func (d *DatabaseConnection) GetDbConnection() (*sql.DB, error) {
	if !d.initialized {
		return nil, errors.New("DatabaseConnection not properly initialized")
	}
	return d.db, nil
}

func (d *DatabaseConnection) Teardown() error {
	if d.db != nil {
		d.db.Close()
	}
	return os.Remove(_DBFILE)
}

func (d *DatabaseConnection) SetUp() {
	dbDir := filepath.Dir(_DBFILE)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.Mkdir(dbDir, os.ModePerm)
		if err != nil {
			log.Info().Err(err).Msg("Creating directory for database files")
			return
		}
	}

	db, err := sql.Open("sqlite3", _DBFILE)
	if err != nil {

		log.Info().Err(err).Msg("Opening sqlite connection")
		return
	}
	if createChatMessagesTable(db) != nil {
		return
	}
	if createEventTable(db) != nil {
		return
	}
	if createChatsTable(db) != nil {
		return
	}
	if createUsersTable(db) != nil {
		return
	}
	d.db = db
	d.initialized = true
}

func (d *DatabaseConnection) IsInitialized() bool {
	return d.initialized
}

func createChatMessagesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ChatMessages (
		id TEXT UNIQUE,
		chat_id TEXT,
		user_id TEXT,
		text TEXT,
		image_base64 TEXT,
		creation_date TEXT,
		"order" INTEGER PRIMARY KEY,
		FOREIGN KEY (chat_id) REFERENCES Chats(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating ChatMessages table: %v", err)
		return err
	}

	return nil
}

func createEventTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating events table: %v", err)
		return err
	}

	return nil
}

func createChatsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Chats (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		creation_date TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating Chats table: %v", err)
		return err
	}

	return nil
}

func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Users (
		user_id TEXT,
		chat_id TEXT,
		FOREIGN KEY (chat_id) REFERENCES Chats(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating Users table: %v", err)
		return err
	}

	return nil
}
