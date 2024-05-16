package models

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
)

var (
	todos = []Todo{
		{
			Id:        1,
			Title:     "task 1",
			Completed: false,
		},
		{
			Id:        2,
			Title:     "task 2",
			Completed: false,
		},
		{
			Id:        3,
			Title:     "task 3",
			Completed: true,
		},
	}
)

type Todo struct {
	Id        int64  `orm:"auto"`
	Title     string `orm:"size(128)"`
	Completed bool
}

func init() {
	orm.RegisterModel(new(Todo))
}

// AddTodo insert a new Todo into database and returns
// last inserted Id on success.
func AddTodo(m *Todo) (id int64, err error) {
	todos = append(todos, *m)

	return int64(len(todos)), nil
}

// GetTodoById retrieves Todo by Id. Returns error if
// Id doesn't exist
func GetTodoById(id int64) (v *Todo, err error) {
	for i := range todos {
		if todos[i].Id == id {
			return &todos[i], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Id : %v , Not Found Todo", id))
}

// GetAllTodo retrieves all Todo matches certain condition. Returns empty list if
// no records exist
func GetAllTodo() (ml []Todo, err error) {
	return todos, nil
}

// UpdateTodo updates Todo by Id and returns error if
// the record to be updated doesn't exist
func UpdateTodoById(m *Todo) (err error) {
	for i := range todos {
		if todos[i].Id == m.Id {
			todos[i] = *m
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Id : %v , Not Found Todo", m.Id))
}

// DeleteTodo deletes Todo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTodo(id int64) (err error) {
	for i := range todos {
		if todos[i].Id == id {
			todos = append(todos[:i], todos[i+1:]...)

			return nil
		}
	}

	return errors.New(fmt.Sprintf("Id : %v , Not Found Todo", id))
}
