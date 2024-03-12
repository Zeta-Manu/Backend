package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Zeta-Manu/Backend/internal/config"
)

type Database struct {
	Conn *sql.DB
}

// NewDatabase creates a new MySQL database connection.
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Check if the database connection is alive
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

// InitializeDatabase initializes and returns a new database connection.
func InitializeDatabase(dbConfig config.DatabaseConfig) (*Database, error) {
	dbDataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
	return NewDatabase(dbDataSourceName)
}

// Close closes the database connection.
func (db *Database) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}
	return nil
}

func (db *Database) CreateTables() error{
    // UserInformation
    _, err =db.Conn.Exec("CREATE TABLE UserInformation (
        email VARCHAR(100) NOT NULL PRIMARY KEY,
        password VARCHAR(255) NOT NULL,
        profile_picture_id INT
    )"
    )
    if err!=nil{
        log.Fatal(error)
    }
    // User
    _, err =db.Conn.Exec("CREATE TABLE User
    (
        uid INT AUTO_INCREMENT PRIMARY KEY,
        email VARCHAR(100) NOT NULL UNIQUE,
        username VARCHAR(50) NOT NULL,
        FOREIGN KEY (email) REFERENCES UserInformation(email)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //UserActivity
    _, err = db.Conn.Exec("CREATE TABLE UserActivity (
        uid INT NOT NULL,
        last_lid INT,
        last_timestamp TIMESTAMP,
        FOREIGN KEY (uid) REFERENCES User(uid)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //UserStat
    _,err=db.Conn.Exec("CREATE TABLE UserStat (
        uid INT NOT NULL,
        history JSON,
        confident JSON,
        FOREIGN KEY (uid) REFERENCES User(uid)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //VideoRecord
    _,err = db.Conn.Exec("CREATE TABLE VideoRecord (
        vid INT AUTO_INCREMENT PRIMARY KEY,
        handsign_id INT,
        user_id INT,
        record_time TIMESTAMP,
        S3_filename VARCHAR(255),
        FOREIGN KEY (handsign_id) REFERENCES HandSign(id),
        FOREIGN KEY (user_id) REFERENCES User(uid)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //Learning
    _,err=db.Conn.Exec("CREATE TABLE Learning (
        lid INT AUTO_INCREMENT PRIMARY KEY,
        lesson_name VARCHAR(255) NOT NULL
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //HandSign
    _,err = db.Conn.Exec("CREATE TABLE HandSign (
        id INT AUTO_INCREMENT PRIMARY KEY,
        handsign VARCHAR(255) NOT NULL
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //Translation
    _,err = db.Conn.Exec("CREATE TABLE Translation (
        id INT AUTO_INCREMENT PRIMARY KEY,
        handsign_id INT,
        language_id INT,
        en_meaning VARCHAR(255) NOT NULL,
        text TEXT,
        FOREIGN KEY (handsign_id) REFERENCES HandSign(id),
        FOREIGN KEY (language_id) REFERENCES Language(id)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //LessonHandSign
    _,err = db.Conn.Exec("CREATE TABLE LessonHandSign (
        id INT AUTO_INCREMENT PRIMARY KEY,
        lesson_id INT,
        handsign_id INT,
        FOREIGN KEY (lesson_id) REFERENCES Learning(lid),
        FOREIGN KEY (handsign_id) REFERENCES HandSign(id)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //Language
    _,err = db.Conn.Exec("CREATE TABLE Language (
        id INT AUTO_INCREMENT PRIMARY KEY,
        language_id VARCHAR(10) NOT NULL UNIQUE
    )")
    if err!=nil{
        log.Fatal(error)
    }
    //LessonTranslation
    _err=db.Conn.Exec("CREATE TABLE LessonTranslation (
        id INT AUTO_INCREMENT PRIMARY KEY,
        lesson_id INT,
        language_id INT,
        translation_text TEXT,
        FOREIGN KEY (lesson_id) REFERENCES Learning(lid),
        FOREIGN KEY (language_id) REFERENCES Language(id)
    )")
    if err!=nil{
        log.Fatal(error)
    }
    fmt.Println("Tables Created")
}