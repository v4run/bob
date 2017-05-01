package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"github.com/v4run/bob/app"
)

// Action defines the different type of actions for the server
type Action string

func (a Action) asPath() string {
	return fmt.Sprintf("/%s", a)
}

func (a Action) String() string {
	return string(a)
}

const (
	// STOPSERVICE is used to kill all jobs and stop the background service
	STOPSERVICE Action = "close"
	// LISTJOBS is used to list all the currently active jobs
	LISTJOBS Action = "list"
	// STOPJOB is used to stop a particular job
	STOPJOB Action = "kill"
	// NEWJOB is used to add a new job
	NEWJOB Action = "new"
	// STATUS is used to return the status of server
	STATUS Action = "status"
	// ROOT defines the root path for the server
	ROOT Action = ""
	// ATTACH is used to attach to a running job
	ATTACH Action = "attach"
)

// S defines the interface for the server
type S interface {
	Serve() error
}

type server struct {
	sync.RWMutex
	*http.ServeMux
	upgrader      websocket.Upgrader
	jobs          map[int]*Job
	registerChan  chan *Job
	deleteChan    chan int
	commEventChan chan comm
	jID           int
}

// Job defines a watch job
type Job struct {
	w   *app.Watcher
	Dir string `json:"dir"`
	ID  int    `json:"id"`
}

func newJob(dir, appName string) *Job {
	return &Job{
		w:   app.NewWatcher([]string{dir}, appName),
		Dir: dir,
	}
}

type comm struct {
	jID    int
	id     string
	op     string
	reader io.Reader
	writer io.Writer
}

// New creates and returns a new instance of server
func New() S {
	s := &server{
		ServeMux: &http.ServeMux{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		jobs:          make(map[int]*Job),
		registerChan:  make(chan *Job),
		deleteChan:    make(chan int),
		commEventChan: make(chan comm),
	}
	go s.jobWatcher()
	return s
}

func (s *server) jobList() []*Job {
	s.RLock()
	defer s.RUnlock()
	m := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		m = append(m, j)
	}
	return m
}

func (s *server) jobWatcher() {
	for {
		select {
		case j := <-s.registerChan:
			j.ID = s.jID
			s.Lock()
			s.jobs[j.ID] = j
			s.Unlock()
			s.jID++
			go j.w.Watch()
		case id := <-s.deleteChan:
			jww.DEBUG.Println("Deleting job", s.jobs[id])
			if _, present := s.jobs[id]; present {
				s.Lock()
				s.jobs[id].w.EventChan <- app.STOPWATCHING
				delete(s.jobs, id)
				s.Unlock()
			}
		case iop := <-s.commEventChan:
			switch iop.op {
			case "attach":
				if _, present := s.jobs[iop.jID]; present {
					jww.DEBUG.Println("Attaching to job", s.jobs[iop.jID])
					fmt.Println("Lock start")
					s.Lock()
					fmt.Println("Lock success")
					s.jobs[iop.jID].w.R.SetOut(iop.writer, iop.id)
					s.jobs[iop.jID].w.R.SetErr(iop.writer, iop.id)
					fmt.Println("Unlock start")
					s.Unlock()
					fmt.Println("Unlock success")
				}
			case "detach":
				if _, present := s.jobs[iop.jID]; present {
					jww.DEBUG.Println("Detaching from job", s.jobs[iop.jID])
					fmt.Println("Lock start")
					s.Lock()
					fmt.Println("Lock success")
					s.jobs[iop.jID].w.R.UnsetOut(iop.id)
					s.jobs[iop.jID].w.R.UnsetErr(iop.id)
					fmt.Println("Unlock start")
					s.Unlock()
					fmt.Println("Unlock success")
				}
			}
		}
	}
}

func (s *server) stop() {
	s.RLock()
	for _, j := range s.jobs {
		j.w.EventChan <- app.STOPWATCHING
	}
	s.RUnlock()
	os.Exit(0)
}

func (s *server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	_ = conn
}

func (s *server) Serve() error {
	s.AddRoutes()
	// s.ServeWS()
	ser := http.Server{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Handler:      s,
	}
	return ser.ListenAndServe()
}

func (s *server) listJobs(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println(r.Method, "Request to list jobs")
	d, _ := json.Marshal(s.jobList())
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(d)
}

func (s *server) newJob(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println(r.Method, "Request to create a new job. Path: ", r.FormValue("path"))
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	info, err := os.Stat(r.FormValue("path"))
	if err != nil || !info.IsDir() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid value `%s` for `path`", r.FormValue("path"))
		return
	}
	var j *Job
	if r.FormValue("app") != "" {
		j = newJob(r.FormValue("path"), r.FormValue("app"))
	} else {
		j = newJob(r.FormValue("path"), filepath.Base(r.FormValue("path")))
	}
	s.registerChan <- j
}

func (s *server) stopJob(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println(r.Method, "Request to stop job", r.FormValue("id"))
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if id, err := strconv.Atoi(r.FormValue("id")); err == nil {
		s.deleteChan <- id
	} else {
		jww.ERROR.Printf("Invalid job id %v. Error: %v", r.FormValue("id"), err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}

func (s *server) stopService(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println(r.Method, STATUS.asPath())
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	go s.stop()
}

func (s *server) status(w http.ResponseWriter, r *http.Request) {
	jww.DEBUG.Println(r.Method, STATUS.asPath())
	w.Write([]byte("OK"))
}

func (s *server) attach(w http.ResponseWriter, r *http.Request) {
	if jID, err := strconv.Atoi(r.FormValue("id")); err == nil {
		c, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			jww.ERROR.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer c.Close()
		pr, pw := io.Pipe()
		id := xid.New().String()
		s.commEventChan <- comm{
			op:     "attach",
			id:     id,
			jID:    jID,
			writer: pw,
		}

		for {
			line, _, err := bufio.NewReader(pr).ReadLine()
			err = c.WriteMessage(websocket.TextMessage, line)
			if err != nil {
				s.commEventChan <- comm{
					op:  "detach",
					id:  id,
					jID: jID,
				}
				break
			}
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (s *server) AddRoutes() {
	jww.DEBUG.Println("Starting web server.")
	s.Handle(ROOT.asPath(), http.FileServer(http.Dir("server/static")))
	s.HandleFunc(STATUS.asPath(), s.status)
	s.HandleFunc(STOPSERVICE.asPath(), s.stopService)
	s.HandleFunc(STOPJOB.asPath(), s.stopJob)
	s.HandleFunc(LISTJOBS.asPath(), s.listJobs)
	s.HandleFunc(NEWJOB.asPath(), s.newJob)
	s.HandleFunc(ATTACH.asPath(), s.attach)
}

func (s *server) ServeWS() {
	// s.Handle("/ws", nil)
}
