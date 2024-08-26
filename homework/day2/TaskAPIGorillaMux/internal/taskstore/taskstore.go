package taskstore

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// TaskStore is a simple in-memory database of tasks;
type TaskStore struct {
	sync.Mutex
	db *sql.DB
	//tasks  map[int]Task
	 //nextId int
}

// TaskStore constructor
func New() *TaskStore {
	os.Remove("./tasks.db")
	ts := &TaskStore{}
	// ts.tasks = make(map[int]Task)
	// ts.nextId = 1
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		fmt.Errorf("Failed open db")
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS task (id INTEGER NOT NULL PRIMARY KEY , text TEXT, tags TEXT,  due TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Errorf("Failed create table")
	}

	ts.db = db
	return ts
}

// CreateTask create a new task in the store
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	// task := Task{
	// 	Id:   ts.nextId,
	// 	Text: text,
	// 	Due:  due}
	// task.Tags = make([]string, len(tags))
	// copy(task.Tags, tags)
	// // Сохранили task в TaskStore
	// ts.tasks[ts.nextId] = task
	// ts.nextId++
	tagsString := strings.Join(tags, ",")
	res, err := ts.db.Exec(`INSERT INTO task (text, tags, due) VALUES (?, ?, ?)`, text, tagsString, due)
	if err != nil {
		fmt.Errorf("Failed to insert task: %v", err)
	   }
	  
	   id, err := res.LastInsertId()
	   if err != nil {
		fmt.Errorf("Failed to get id: %v", err)
	   }
	  
	   return int(id)
}

// GetTask retrieves the task from taskstore by given id
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	// t, ok := ts.tasks[id]
	// if ok {
	// 	return t, nil
	// } else {
	// 	return Task{}, fmt.Errorf("task with id=%d not found", id)
	// }
	
	var task Task
	var tagsString string
	var dueString string
	err := ts.db.QueryRow(`SELECT id, text, tags, due FROM task WHERE id = ?`, id).Scan(&task.Id, &task.Text, &tagsString, &dueString)
	if err != nil {
	  return Task{}, fmt.Errorf("task with id=%d not found", id)
	 }
	 //Приводим время в формат time.Time
	task.Due, err = time.Parse("2006-01-02 15:04:05 Z0700 MST", dueString)
	//Приводим теги в формат списка
	task.Tags = strings.Split(tagsString, ",")
	return task, nil
}

// GetAllTask retrieves all task from taskstore, in arbitrary order
func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	// allTasks := make([]Task, 0, len(ts.tasks))
	// for _, task := range ts.tasks {
	// 	allTasks = append(allTasks, task)
	// }
	rows, err := ts.db.Query(`SELECT id, text, tags, due FROM task`)
	if err != nil {
		fmt.Errorf("Failed to get tasks: %v", err)
	}
	defer rows.Close()
   
	//Заводим массив для всех найденных тасок
	var allTasks []Task
	for rows.Next() {
	 var task Task
	 var tagsString string
	 var dueString string
	 err := rows.Scan(&task.Id, &task.Text, &tagsString, &dueString)
	 if err != nil {
		fmt.Errorf("Failed to scan task: %v", err)
	 }
	 // Приводим время в формат time.Time
	 task.Due, err = time.Parse("2006-01-02 15:04:05 Z0700 MST", dueString)
	 // Приводим теги в формат списка
	 task.Tags = strings.Split(tagsString, ",")

	 allTasks = append(allTasks, task)
	}
   
	return allTasks

}

// Возвращаем таски с заданным тэгом
func (ts *TaskStore) GetTaskbyTag(tag string) ([]Task, error) {
	ts.Lock()
	defer ts.Unlock()

	// allTasksWithTag := make([]Task, 0)
	// for _, task := range ts.tasks {
	// 	for _, t := range task.Tags {
	// 		if t == tag {
	// 			allTasksWithTag = append(allTasksWithTag, task)
	// 		}
	// 		if len(allTasksWithTag) == 0 {
	// 			return nil, fmt.Errorf("tasks with tag=%s not found", tag)
	// 		}
	// 	}
	// }
	rows, err := ts.db.Query(`SELECT id, text, tags, due FROM task WHERE tags LIKE ?`, "%"+tag+"%")
	if err != nil {
		fmt.Errorf("Failed to get tasks: %v", err)
	}
	defer rows.Close()
   
	// Заводим массив для всех найденных тасок по тегам
	var allTasksWithTag []Task
	for rows.Next() {
	 var task Task
	 var tagsString string
	 var dueString string
	 err := rows.Scan(&task.Id, &task.Text, &tagsString, &dueString)
	 if err != nil {
		fmt.Errorf("Failed to scan task: %v", err)
	 }
	 // Приводим время в формат time.Time
	 task.Due, err = time.Parse("2006-01-02 15:04:05 Z0700 MST", dueString)
	 // Приводим теги в формат списка
	 task.Tags = strings.Split(tagsString, ",")

	 allTasksWithTag = append(allTasksWithTag, task)

	if len(allTasksWithTag) == 0 {
		return nil, fmt.Errorf("tasks with tag=%s not found", tag)
		}
	}
	return allTasksWithTag, nil
}

// Возвращаем таски с заданной датой (без учета времени)
func (ts *TaskStore) GetTaskByDate(date time.Time) ([]Task, error) {
	ts.Lock()
	defer ts.Unlock()

	// var tasksOnDate []Task
	// for _, task := range ts.tasks {
	// 	if task.Due.Year() == date.Year() && task.Due.Month() == date.Month() && task.Due.Day() == date.Day() {
	// 		tasksOnDate = append(tasksOnDate, task)
	// 	}
	// }

	// Приводим time.Time в sting для SQL-запроса
	dateStr := date.Format("2006-01-02")
	rows, err := ts.db.Query(`SELECT id, text, tags, due FROM task WHERE due LIKE ?`, "%"+dateStr+"%")
	if err != nil {
		fmt.Errorf("Failed to get tasks: %v", err)
	}
	defer rows.Close()
   
	// Заводим массив для всех найденных тасок по тегам
	var tasksOnDate []Task
	for rows.Next() {
	 var task Task
	 var tagsString string
	 var dueString string
	 err := rows.Scan(&task.Id, &task.Text, &tagsString, &dueString)
	 if err != nil {
		fmt.Errorf("Failed to scan task: %v", err)
	 }
	 // Приводим время в формат time.Time
	 task.Due, err = time.Parse("2006-01-02 15:04:05 Z0700 MST", dueString)
	 // Приводим теги в формат списка
	 task.Tags = strings.Split(tagsString, ",")

	 tasksOnDate = append(tasksOnDate, task)

	if len(tasksOnDate) == 0 {
		return nil, fmt.Errorf("tasks on date=%v not found", date)
		}
	}

	return tasksOnDate, nil
}

// DeleteAllTasks deletes all tasks in the taskstore
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	_, err := ts.db.Exec("DELETE from task")
	if err != nil {
		fmt.Errorf("Failed to erase task db: %v", err)
	}
	return nil
}

// DeleteTask deletes the task from taskstore by given id. If no such id exists, return Error
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()
   
	_, err := ts.db.Exec(`DELETE FROM task WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("task with id=%d not found", id)
	}
	return nil
}
