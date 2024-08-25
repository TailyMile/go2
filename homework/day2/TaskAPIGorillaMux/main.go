/*
// Пример REST сервера с несколькими маршрутами(библиотеку GorillaMux)

// POST   /task/              :  создаёт задачу и возвращает её ID
// GET    /task/<taskid>      :  возвращает одну задачу по её ID
// GET    /task/              :  возвращает все задачи
// DELETE /task/<taskid>      :  удаляет задачу по ID
// DELETE /task/              :  удаляет все задачи
// GET    /tag/<tagname>      :  возвращает список задач с заданным тегом
// GET    /due/<yy>/<mm>/<dd> :  возвращает список задач, запланированных на указанную дату

*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/vadshi/go2/TaskStoreAPI/internal/taskstore"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func (ts *taskServer) taskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if len(id) == 0 {
		if r.Method == http.MethodPost {
			ts.createTaskHandler(w, r)
		} else if r.Method == http.MethodGet {
			ts.getAllTaskHandler(w, r)
		} else if r.Method == http.MethodDelete {
			ts.deleteAllTaskHandler(w, r)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, POST, DELETE at '/task', got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}

	} else {
		// Приводим id к int
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodGet {
			ts.getTaskHandler(w, r, idInt)
		} else if r.Method == http.MethodDelete {
			ts.deleteTaskHandler(w, r, idInt)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE at '/task/<id>', got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling task create at %s\n", r.URL.Path)

	// Структура для создания задачи
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	// Для ответа в виде JSON
	type ResponseId struct {
		Id int `json:"id"`
	}

	// JSON в качестве Content-Type
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	// Обработка тела запроса
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Создаем новую задачу в хранилище и получаем ее <id>
	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)

	// Создаем json для ответа
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (ts *taskServer) getAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get all tasks at %s\n", r.URL.Path)

	allTasks := ts.store.GetAllTasks()

	js, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("Handling get task at %s\n", r.URL.Path)

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // код ошибки 404
		return
	}
	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request, id int) {
	log.Printf("Handling delete task at %s\n", r.URL.Path)

	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // код ошибки 404
		return
	}
}

func (ts *taskServer) deleteAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling delete all tasks at %s\n", r.URL.Path)

	ts.store.DeleteAllTasks()
}

func (ts *taskServer) getTaskbyTagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get task by tag at %s\n", r.URL.Path)

	vars := mux.Vars(r)
	tag := vars["tag"]
	allTasksWithTag, err := ts.store.GetTaskbyTag(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // код ошибки 404
		return
	}

	js, err := json.Marshal(allTasksWithTag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // код ошибки 500
		return
	}

	// Обязательно вносим изменения в Header до вызова метода Write()!
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTaskByDateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get task by date at %s\n", r.URL.Path)

	// Получаем дату из URL
	vars := mux.Vars(r)
	dateStr := fmt.Sprintf("%s/%s/%s", vars["y"], vars["m"], vars["d"])

	date, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid date format, expected /due/yyyy/mm/dd, got %v", dateStr), http.StatusBadRequest)
		return
	}

	// Получаем задачи по дате
	allTasksByDate, err := ts.store.GetTaskByDate(date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Формируем JSON-ответ
	js, err := json.Marshal(allTasksByDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	// mux := http.NewServeMux()
	r := mux.NewRouter()
	server := NewTaskServer()

	// Added routing for "/task/" URL
	r.HandleFunc("/task/", server.taskHandler).Methods("GET", "POST", "DELETE")
	// Added routing for "/tag/" URL
	r.HandleFunc("/task/{id}", server.taskHandler).Methods("GET", "DELETE")
	r.HandleFunc("/tag/{tag}", server.getTaskbyTagHandler).Methods("GET")
	// Added routing for "/due" URL
	r.HandleFunc("/due/{y}/{m}/{d}", server.getTaskByDateHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
