package main

import (
	"fmt"
	"math/rand"
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

// Simulate a more intensive job with random processing time.
func longDoubler(job Job) Job {
	processTime := time.Duration(rand.Intn(500)) * time.Millisecond // Random delay between 0-500ms
	time.Sleep(processTime)
	return Job{data: job.data, result: job.data * 2}
}

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
			// Simulate a longer job
			results <- longDoubler(job)
		}
	}
}

func main() {
	jobCount := 1000 // Increased number of jobs for stress testing
	jobs := make(chan Job, jobCount)
	results := make(chan Job, jobCount)
	control := make(chan ControlMsg)

	// Launch multiple workers (stress test concurrency).
	workerCount := 5
	for w := 0; w < workerCount; w++ {
		go doubler(jobs, results, control)
	}

	// Send jobs to the jobs channel.
	for i := 0; i < jobCount; i++ {
		jobs <- Job{i, 0}
	}
	
	resultCount := 0
	sum := 0

	timeout := time.Duration(2*jobCount) * time.Millisecond // Adjust timeout based on workload
	start := time.Now()

	for {
		select {
		case result := <-results:
			fmt.Println(result)
			sum += result.result
			resultCount++

			// Check if all jobs are processed.
			if resultCount == jobCount {
				fmt.Printf("All jobs processed. Total time: %v\n", time.Since(start))
				for i := 0; i < workerCount; i++ {
					control <- DoExit
					<-control // Wait for goroutine to confirm exit.
				}
				fmt.Println("Final sum is:", sum)
				fmt.Println("Program exit")
				return
			}

		case <-time.After(timeout):
			fmt.Println("Timed out, forcing exit")
			for i := 0; i < workerCount; i++ {
				control <- DoExit
				<-control // Wait for goroutine to confirm exit.
			}
			fmt.Println("Final sum at timeout:", sum)
			fmt.Println("Program exit")
			return
		}
	}
}
