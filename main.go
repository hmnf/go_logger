package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/hmnf/go_logger/logger"
	"github.com/hmnf/go_logger/store"
)

type Repository struct {
	logger *logger.Logger
}

func NewRepository(logger *logger.Logger) *Repository {
	return &Repository{
		logger: logger,
	}
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var userStorage = store.NewStorage[User]()

func (repo *Repository) GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]

	user, err := userStorage.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(bytes)
}

func (repo *Repository) StoreUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repo.logger.WritePut(name, string(bodyBytes))

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userStorage.Put(name, user)
}

func (repo *Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	repo.logger.WriteDelete(name)
	userStorage.Delete(name)
}

func Overload(w http.ResponseWriter, r *http.Request) {
	wg := &sync.WaitGroup{}
	names := []string{
		"Joe",
		"Emily",
		"Nick",
		"Nensi",
		"Lion",
		"Mike",
		"Anna",
		"Helen",
		"Melony",
		"Kate",
	}
	wg.Add(len(names))
	client := &http.Client{}
	for i := 0; i < len(names); i++ {
		go func(i int) {
			user := User{
				Name: names[i],
				Age:  100,
			}
			userBytes, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
				return
			}
			request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/%s", names[i]), bytes.NewBuffer(userBytes))
			if err != nil {
				log.Println(err)
			}
			resp, err := client.Do(request)
			if err != nil {
				if resp != nil {
					log.Println(resp.Status)
				}
				log.Println(err)
			} else {
				log.Println(resp.Status)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
}

func main() {
	logger, err := logger.NewLogger("test.log")
	if err != nil {
		log.Fatal(err)
	}
	logger.Run()

	// events, errors := logger.ReadEvents()

	// for event := range events {
	//     store.Put(event.key, event.Value)
	// }
	repo := NewRepository(logger)

	router := mux.NewRouter()
	router.HandleFunc("/{name}", repo.GetValue).Methods("GET")
	router.HandleFunc("/{name}", repo.StoreUser).Methods("POST")
	router.HandleFunc("/{name}", repo.DeleteUser).Methods("DELETE")

	router.HandleFunc("/v1/overload", Overload).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// wg := sync.WaitGroup{}
// wgGet := sync.WaitGroup{}
// wgDel := sync.WaitGroup{}

// wg.Add(10)
// wgGet.Add(10)
// wgDel.Add(10)
// for i := 0; i < 10; i++ {
// 	fmt.Println(i)
// 	go func(i int) {
// 		<-time.After(3 * time.Second)
// 		key := fmt.Sprintf("key_%d", i)
// 		value := fmt.Sprintf("value_%d", i)
// 		s.Put(key, value)
// 		wg.Done()
// 	}(i)

// 	go func(i int) {
// 		key := fmt.Sprintf("key_%d", i)
// 		s.Get(key)
// 		wgGet.Done()
// 	}(i)

// 	go func(i int) {
// 		key := fmt.Sprintf("key_%d", i)
// 		s.Delete(key)
// 		wgDel.Done()
// 	}(i)
// }

// wg.Wait()
// wgGet.Wait()
// wgDel.Wait()
// fmt.Println("Finish!")

// func main() {

// ch := make(chan string)

// go func() {
// 	<-time.After(2 * time.Second)
// 	ch <- "hello"
// }()

// fmt.Println("Waiting...")
// fmt.Println(<-ch)

// ch := make(chan string, 3)

// ch <- "Hello"
// ch <- "foo"
// ch <- "bar"
// close(ch)

// for m := range ch {
// 	fmt.Println(m)
// }

// msg, ok := <-ch
// fmt.Printf("%q, %v\n", msg, ok)

// ch1 := make(chan string, 2)
// ch2 := make(chan string, 2)
// ch3 := make(chan string, 2)

// go func() {
// 	<-time.After(2 * time.Second)
// 	ch1 <- "channel 1"
// }()

// go func() {
// 	<-time.After(1 * time.Second)
// 	ch2 <- "channel 2"
// }()

// go func() {
// 	<-time.After(2 * time.Second)
// 	ch3 <- "channel 3"
// }()

// select {
// case <-ch1:
// 	fmt.Println("Channel 1")
// case x := <-ch2:
// 	fmt.Println(x)
// case <-ch3:
// 	fmt.Println("Channel 3")
// case <-time.After(10 * time.Second):
// 	fmt.Println("Nothing")
// }
// }
