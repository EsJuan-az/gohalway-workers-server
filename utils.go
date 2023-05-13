package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Job struct {
	Name string
	Delay time.Duration
	Number int 
}
type Worker struct{
	Id int
	JobQueue chan Job
	WorkerPool chan chan Job
	QuitChan chan bool
}
func (w *Worker)Start(){
	go func (){
		for {
			w.WorkerPool <- w.JobQueue
			select{
			case job := <-w.JobQueue:
				fmt.Printf("Worker %d has started with %d\n", w.Id, job.Number)
				res := Fibo(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker %d has finished with %d: %d gotten as a result\n", w.Id, job.Number, res)
			case <- w.QuitChan:
				fmt.Printf("Worker %d has stopped\n", w.Id)
			}
		}
	}()
}
func (w *Worker)Stop(){
	go func(){
		w.QuitChan <- true
	}()
}
func NewWorker(id int, workerpool chan chan Job) *Worker{
	return &Worker{
		Id: id,
		JobQueue: make(chan Job),
		WorkerPool: workerpool,
		QuitChan: make(chan bool),
	}
}
type Dispatcher struct{
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue chan Job
}
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher{
	return &Dispatcher{
		JobQueue: jobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: make(chan chan Job, maxWorkers),
	}
}
func (d *Dispatcher) Dispatch(){
	for{
		select{
		case job := <-d.JobQueue:
			go func(){
				workerJobQueue := <- d.WorkerPool
				workerJobQueue <- job
			}()
		}
	}
}
func (d *Dispatcher) Run(){
	for i := 0; i < d.MaxWorkers; i ++{
		w := NewWorker(i, d.WorkerPool)
		w.Start()
	}
	go d.Dispatch()
}
func Fibo(n int)int {
	if n <= 2{
		return 1
	}
	return Fibo(n-1) + Fibo(n-2)
}
func RequestHandler(jobQueue chan Job)http.HandlerFunc{
	return func(w http.ResponseWriter, request *http.Request){
		if request.Method != "POST"{
			w.Header().Set("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		delay, err1 := time.ParseDuration(request.FormValue("delay"))
		value, err2 := strconv.Atoi(request.FormValue("number"))
		name := request.FormValue("name")
		if err1 != nil || err2 != nil || name == ""{
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		job := Job{
			Name: name, 
			Delay: delay,
			Number: value,
		}
		jobQueue <- job
		w.WriteHeader(http.StatusCreated)
	}
}