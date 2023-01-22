package main

import (
	"log"
	"net/http"

	"github.com/hmnf/go_logger/logger"
	"github.com/hmnf/go_logger/store"
)

func main() {
	s, err := store.NewStorageService()
	s.Restore()
	if err != nil {
		log.Fatal(err)
		return
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
	l, err := logger.NewLogger()

	if err != nil {
		log.Fatal(err)
		return
	}

	l.Run()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		s.Get("Joe")
		w.Write([]byte("Success"))
	})
	http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		s.Put("Joe", "User Name")
		l.WritePut("Put message")
		w.Write([]byte("Success"))
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		s.Delete("Joe")
		l.WriteDelete("Delete message")
		w.Write([]byte("Success"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
