package engine

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	"web-service-kubernetes/datastore"
	"web-service-kubernetes/servicelog"
	"web-service-kubernetes/utility"
)

// Get all collection
func GetCollection(w http.ResponseWriter, r *http.Request) {
	// Send the all collection elements
	todos, err := datastore.GetCollection()
	if err != nil {
		logger := servicelog.GetInstance()
		logger.Println(time.Now().UTC(), "Get Collection failed")
		// send error
		utility.EncodeToJsonError(w)
		return
	}
	utility.EncodeToJsonWithBody(w, todos)
}

// Get an element of the collection
func GetElement(w http.ResponseWriter, r *http.Request) {
	logger := servicelog.GetInstance()
	var todo *datastore.TodoElement
	// taking the id
	indexString := r.URL.Path[len("/todos/"):]
	if todoId, err := strconv.Atoi(indexString); err == nil {
		todo, err = datastore.Get(todoId)
		if err != nil {
			// Not found
			if err == sql.ErrNoRows {
				utility.EncodeToJsonNotFound(w)
				return
			}
			logger.Println(time.Now().UTC(), "GetElement failed")
			// send error
			utility.EncodeToJsonError(w)
			return
		}
		// send the element requested
		utility.EncodeToJsonWithBody(w, todo)
	} else {
		logger.Println(time.Now().UTC(), "GetElement failed")
		// send error
		utility.EncodeToJsonError(w)
	}
}

// Create a new todoItem element
func CreateElement(w http.ResponseWriter, r *http.Request) {
	var newTodo datastore.TodoElement
	var err error
	var id int64
	logger := servicelog.GetInstance()
	if newTodo, err = utility.MarshallJsonAndResponse(w, r); err != nil {
		logger.Println(time.Now().UTC(), "CreateElement failed")
		// returns error
		utility.EncodeToJsonError(w)
		return
	}
	if id, err = datastore.Put(newTodo); err != nil {
		logger.Println(time.Now().UTC(), "CreateElement failed")
		// send error
		utility.EncodeToJsonError(w)
		return
	}
	// Send created with new resource id
	utility.EncodeToJson(w, r, id)
}

// Update an element
func UpdateElement(w http.ResponseWriter, r *http.Request) {
	var updatedTodo datastore.TodoElement
	var err error
	var updatedid int64
	logger := servicelog.GetInstance()
	if updatedTodo, err = utility.MarshallJsonAndResponse(w, r); err != nil {
		// returns error
		utility.EncodeToJsonError(w)
		return
	}

	if updatedid, err = datastore.Update(updatedTodo); err != nil {
		// Not found
		if err == sql.ErrNoRows {
			utility.EncodeToJsonNotFound(w)
			return
		}
		logger.Println(time.Now().UTC(), "UpdateElement failed")
		// send error
		utility.EncodeToJsonError(w)
		return
	}
	// Send created with new resource id
	utility.EncodeToJson(w, r, updatedid)
}

// Delete all collection
func DeleteCollection(w http.ResponseWriter, r *http.Request) {
	datastore.DeleteCollection()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Delete an element of the collection
func DeleteElement(w http.ResponseWriter, r *http.Request) {
	// taking the id
	indexString := r.URL.Path[len("/todos/"):]
	logger := servicelog.GetInstance()
	if todoId, err := strconv.Atoi(indexString); err == nil {
		if err := datastore.DeleteElement(todoId); err != nil {
			// Not found
			if err == sql.ErrNoRows {
				utility.EncodeToJsonNotFound(w)
				return
			}
			logger.Println(time.Now().UTC(), "DeleteElement failed")
			// send error
			utility.EncodeToJsonError(w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		utility.EncodeToJsonError(w)
	}
}
