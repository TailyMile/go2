package taskstore

import (
	"testing"
	"time"
)

func TestCreateAndGet(t *testing.T){
	// Создаем хранилище и задачу
		ts := New()
		id := ts.CreateTask("Check", nil, time.Now())

		task, err := ts.GetTask(id)
		if err != nil{
			t.Fatal(err)
		}

		if task.Id != id  {
			t.Errorf("got task with Id=%d but not id=%d", task.Id, id)
		}
		if task.Text != "Check"{
			t.Errorf("got task with Text=%v, want %v", task.Text, "Check")
		}
		
		_, err = ts.GetTask(id + 1)
		if err == nil{
			t.Errorf("got nil, want error.")
		}

		ts.CreateTask("Second check", nil, time.Now())

		allTasks := ts.GetAllTasks()
		if len(allTasks) != 2{
			t.Errorf("got len(allTasks)=%d, want 2", len(allTasks))
		}


	}