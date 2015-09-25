package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	//"html"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func TodoIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(todos)
}

func TodoShow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	i, err := strconv.Atoi(ps.ByName("todoId"))
	if err != nil {
		fmt.Fprintln(w, "wrong id: "+ps.ByName("todoId"))
		return
	}
	fmt.Fprintln(w, "Todo show: "+strconv.Itoa(i))

	var t *Todo
	for _, v := range todos {
		if v.Id == i {
			t = &v
			break
		}
	}

	if t != nil {
		json.NewEncoder(w).Encode(t)
	} else {
		fmt.Fprintln(w, "Not found")
	}
}

func TodoCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	t := AddTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

type Todos []Todo

func AddTodo(t Todo) Todo {
	index++
	t.Id = index
	todos = append(todos, t)
	return t
}

var (
	todos = make(Todos, 0)
	index = 0
)

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/todos", TodoIndex)
	router.GET("/todos/:todoId", TodoShow)
	router.POST("/todos/", TodoCreate)

	log.Fatal(http.ListenAndServe(":8080", router))
}
