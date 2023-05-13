package main

import (
	"fmt"
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
func Fibo(n int)int {
	if n <= 2{
		return 1
	}
	return Fibo(n-1) + Fibo(n-2)
}