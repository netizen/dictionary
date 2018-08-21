package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func constructDictionary() {
	file, err := ioutil.ReadFile("dictionary.json")
	check(err)

	var m map[string]string
	err = json.Unmarshal(file, &m)
	check(err)

	for k, v := range m {
		w := Word{k, v}
		Dictionary = append(Dictionary, w)
	}
	sort.Slice(Dictionary, func(i, j int) bool {
		return Dictionary[i].Term < Dictionary[j].Term
	})
	log.Printf("Dictionary loaded with %d words... ", len(Dictionary))
}

type Word struct {
	Term       string
	Definition string
}

var Dictionary = []Word{}

type Worker struct {
	words []Word
	ch    chan *Word
	name  string
}

func NewWorker(words []Word, ch chan *Word, name string) *Worker {
	return &Worker{words: words, ch: ch, name: name}
}

func (w *Worker) Find(term string) {
	for i := range w.words {
		word := &w.words[i]
		if strings.ToLower(word.Term) == strings.ToLower(term) {
			log.Printf("Worker %s found a matching entry", w.name)
			log.Println(word)
			w.ch <- word
		}
	}
}

var worker1, worker2, worker3, worker4 *Worker
var ch = make(chan *Word)

func findDefinition(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving a request...")
	term := r.URL.Path
	term = strings.TrimPrefix(term, "/")

	switch {
	case term[0] <= 'f':
		log.Println("Delegating to worker1")
		go worker1.Find(term)
	case term[0] <= 'l':
		log.Println("Delegating to worker2")
		go worker2.Find(term)
	case term[0] <= 'r':
		log.Println("Delegating to worker3")
		go worker3.Find(term)
	default:
		log.Println("Delegating to worker4")
		go worker4.Find(term)
	}

	for {
		select {
		case word := <-ch:
			b, err := json.Marshal(word)
			check(err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(b)
		case <-time.After(100 * time.Millisecond):
			return
		}
	}
}

func initializeWorkers() {
	ch = make(chan *Word)
	var part1, part2, part3 int
	for i := range Dictionary {
		if (strings.ToLower(Dictionary[i].Term))[0] == 'f' {
			part1 = i
		} else if (strings.ToLower(Dictionary[i].Term))[0] == 'l' {
			part2 = i
		} else if (strings.ToLower(Dictionary[i].Term))[0] == 'r' {
			part3 = i
		}
	}

	worker1 = NewWorker(Dictionary[:part1], ch, "#1")
	worker2 = NewWorker(Dictionary[part1:part2], ch, "#2")
	worker3 = NewWorker(Dictionary[part2:part3], ch, "#3")
	worker4 = NewWorker(Dictionary[part3:], ch, "#4")

	log.Printf("Workers initialized %d %d %d %d", len(Dictionary[:part1]), len(Dictionary[part1:part2]), len(Dictionary[part2:part3]), len(Dictionary[part3:]))
}

func main() {
	constructDictionary()
	initializeWorkers()
	http.HandleFunc("/", findDefinition)
	log.Println("Now ready to serve requests")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
