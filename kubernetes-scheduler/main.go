package main

import (
	"fmt"
	"k8s.io/api/core/v1"
)

func main() {
	fmt.Println("Start a scheduler")

	podQueue := make(chan *v1.Pod, 300)
	defer close(podQueue)

	quit := make(chan struct{})
	defer close(quit)

	scheduler := NewScheduler(podQueue, quit)
	scheduler.Run()
}

type Scheduler struct {
}

func (s *Scheduler) Run() {

}

func NewScheduler(podQueue chan *v1.Pod, quit chan struct{}) Scheduler {
	return Scheduler{}
}
