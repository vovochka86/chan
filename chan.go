package main

import (
	"fmt"
	"time"
)

type ControlMsg int

type Job struct {
	data   int
	result int
}

const (
	DoExit ControlMsg = iota
	ExitOk
)

func doubler(jobs, results chan Job, control chan ControlMsg) {
	for {
		select {
		case msg := <-control:
			switch msg {
			case DoExit:
				fmt.Println("Exit goroutine")
				control <- ExitOk
				return
			default:
				panic("Unhandled control message")
			}
		case job := <-jobs:
			results <- Job{data: job.data, result: job.data * 2}
		//	default:
								//			 time.Sleep(50 * time.Millisecond)
		}
	}
}

func main() {
	jobs := make(chan Job, 30)   // Smaller buffer size, fits exactly 30 jobs.
	results := make(chan Job, 30) // Results channel matches jobs count.
	control := make(chan ControlMsg)

	// Start the worker goroutine.
	go doubler(jobs, results, control)

	// Send jobs to the jobs channel.
	for i := 0; i < 30; i++ {
		jobs <- Job{i, 0}
	}
	//close(jobs) // Close the jobs channel as no more jobs will be sent.

	// Process the results or time out after a while.
	for {
		select {
		case result := <-results:
			fmt.Println(result)
			if len(results) == 0 && len(jobs) == 0 {
				// If no more jobs or results, exit.
				fmt.Println("All jobs processed")
				control <- DoExit
				<-control
				fmt.Println("Program exit")
				return
			}
		case <-time.After(500 * time.Millisecond):
			fmt.Println("timed out")
			control <- DoExit
			<-control // Wait for the goroutine to confirm exit.
			fmt.Println("Program exit")
			return
		}
	}
}
