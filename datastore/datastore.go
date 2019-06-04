/* The datastore package gives an mysql implementation to save the todos in the DBMS */
package datastore

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
	"web-service-kubernetes/servicelog"

	_ "github.com/go-sql-driver/mysql"
)

// A todoElement element
type TodoElement struct {
	Id    int    `json:"Id"`
	Topic string `json:"Topic"`
	// sql use 0/1 for boolean
	Completed int    `json:"Completed"`
	Due       string `json:"Due"`
}

// A list of todos element (used to unmarshal a bunch of todos elements from a file
type TodoElements struct {
	Todos []TodoElement
}

// Connect to the mysql database
func connectDatastore() (*sql.DB, error) {
	var credentials string
	var host string
	var database string
	var err error
	logger := servicelog.GetInstance()
	logger.Println("connecting to the datastore")
	if credentials, host, database, err = LoadConfiguration(); err != nil {
		logger.Println(time.Now().UTC(), "Error loading configuration")
		return nil, err
	}
	// Create an sql.DB and check for error
	connstring := credentials + "tcp(" + host + ")/" + database
	logger.Println(time.Now().UTC(), connstring)
	db, err := sql.Open("mysql", connstring)
	if err != nil {
		logger.Println(time.Now().UTC(), "Error connecting to the database")
		return nil, err
	}
	logger.Println("returning from connecting to the datastore")
	return db, nil
}

// Get all the collection (GET on todos)
func GetCollection() ([]TodoElement, error) {
	// Get all the element from the database
	todos := make([]TodoElement, 0)
	var db *sql.DB
	var err error
	var rows *sql.Rows
	logger := servicelog.GetInstance()
	logger.Println("retriving all elements from collection")
	if db, err = connectDatastore(); err != nil {
		logger.Println(time.Now().UTC(), "Error Connecting to datastore in GetCollection")
		return nil, err
	}
	logger.Println("before defferring")
	defer db.Close()
	if rows, err = db.Query("SELECT ID, Topic, Completed, Due FROM ToDo"); err != nil {
		logger.Println(time.Now().UTC(), "Error executing query in GetCollection"+err.Error())
		return nil, err
	}
	logger.Println("after selecting")
	defer rows.Close()
	for rows.Next() {
		elem := new(TodoElement)
		if err := rows.Scan(&elem.Id, &elem.Topic, &elem.Completed, &elem.Due); err != nil {
			logger.Println(time.Now().UTC(), "Error Scanning rows in Get Collection")
			return todos, err
		}
		todos = append(todos, *elem)
	}
	if err := rows.Err(); err != nil {
		return todos, err
	}
	logger.Println("ending retriving all elements from collection")
	return todos, nil
}

func Get(index int) (*TodoElement, error) {
	var db *sql.DB
	var err error
	todo := new(TodoElement)
	logger := servicelog.GetInstance()
	if db, err = connectDatastore(); err != nil {
		logger.Println(time.Now().UTC(), "Error Connecting to datastore in Get")
		return nil, err
	}
	defer db.Close()
	if err = db.QueryRow("SELECT ID, Topic, Completed, Due FROM ToDo WHERE ID=?", index).Scan(&todo.Id, &todo.Topic, &todo.Completed, &todo.Due); err != nil {
		logger.Println(time.Now().UTC(), "Error executing query in Get")
		return nil, err
	}

	return todo, nil
}

// Put an element in the map
func Put(todo TodoElement) (int64, error) {
	var db *sql.DB
	var err error
	var id int64
	var r sql.Result
	logger := servicelog.GetInstance()
	if db, err = connectDatastore(); err != nil {
		logger.Println(time.Now().UTC(), "Error Connecting to datastore in Put")
		return id, err
	}
	defer db.Close()
	r, err = db.Exec("INSERT INTO ToDo(Topic, Completed, Due) VALUES(?, ?, ?)", todo.Topic, todo.Completed, todo.Due)
	if err != nil {
		logger.Println(time.Now().UTC(), "Error executing query in Put")
		return id, err
	}
	id, _ = r.LastInsertId()
	return id, nil
}

// Update an element in the map
func Update(todo TodoElement) (int64, error) {
	var db *sql.DB
	var err error
	var id int64
	var r sql.Result
	logger := servicelog.GetInstance()
	if db, err = connectDatastore(); err != nil {
		logger.Println(time.Now().UTC(), "Error Connecting to datastore in Update")
		return id, err
	}
	defer db.Close()
	r, err = db.Exec("UPDATE ToDo SET Topic=?, Completed=?, Due=? WHERE ID=?", todo.Topic, todo.Completed, todo.Due, todo.Id)
	if err != nil {
		logger.Println(time.Now().UTC(), "Error executing query in Put")
		return id, err
	}
	id, _ = r.LastInsertId()
	return id, nil
}

// Update an element in the map
func DeleteElement(index int) error {
	var db *sql.DB
	var err error
	logger := servicelog.GetInstance()
	if db, err = connectDatastore(); err != nil {
		logger.Println(time.Now().UTC(), "Error Connecting to datastore in DeleteElement")
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ToDo WHERE ID=?", index)
	if err != nil {
		logger.Println(time.Now().UTC(), "Error executing query in DeleteCollection")
		return err
	}
	return nil
}

// Clear the map (delete operation)
func DeleteCollection() error {
	var db *sql.DB
	var err error
	logger := servicelog.GetInstance()
	if db, err = connectDatastore(); err != nil {
		logger.Println(time.Now().UTC(), "Error Connecting to datastore in DeleteCollection")
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM ToDo")
	if err != nil {
		logger.Println(time.Now().UTC(), "Error executing query in DeleteCollection")
		return err
	}
	return nil
}

// Load DBMS credentials and host
func LoadConfiguration() (string, string, string, error) {
	var file *os.File
	var err error
	logger := servicelog.GetInstance()
	if file, err = os.Open("conf"); err != nil {
		logger.Println(time.Now().UTC(), "Error opening configuration file")
		// load default
		return "", "", "", nil
		//return "root:my-secret-pw@", "172.17.0.2:3306", nil
	}
	var Credentials string = ""
	var Host string = ""
	var Database string = ""
	if err != nil {
		Credentials = ""
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	Username, _ := scanner.ReadString(':')
	Username, _ = scanner.ReadString('\n')

	Username = Username[:len(Username)-1]
	Passwd, _ := scanner.ReadString(':')
	Passwd, _ = scanner.ReadString('\n')
	Passwd = Passwd[:len(Passwd)-1]

	if strings.Compare(Username, "") == 0 || strings.Compare(Passwd, "") == 0 {
		Credentials = ""
	} else {
		Credentials = fmt.Sprintf("%s:%s@", Username, Passwd)
	}
	Host, _ = scanner.ReadString(':')
	Host, _ = scanner.ReadString('\n')
	Host = Host[:len(Host)-1]

	Database, _ = scanner.ReadString(':')
	Database, _ = scanner.ReadString('\n')
	Database = Database[:len(Database)-1]

	Port, _ := scanner.ReadString(':')
	Port, _ = scanner.ReadString('\n')
	Port = Port[:len(Port)-1]
	Host = fmt.Sprintf("%s:%s", Host, Port)

	return Credentials, Host, Database, nil
}
