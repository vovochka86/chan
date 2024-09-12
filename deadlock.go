package main

import (
	"fmt"
	"sync"
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

func doubler(jobs, results chan Job, control chan ControlMsg, wg *sync.WaitGroup) {
	defer wg.Done() // Notify when goroutine finishes

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
		case job, ok := <-jobs:
			if !ok { // If jobs channel is closed, return
				return
			}
			results <- Job{data: job.data, result: job.data * 2}
		}
	}
}

func main() {
	numJobs := 5000            // Increase the number of jobs for stress testing
	numWorkers := 10           // Spawn multiple workers for concurrency
	jobs := make(chan Job, 100) // Buffered channel for job handling
	results := make(chan Job, 100)
	control := make(chan ControlMsg)

	var wg sync.WaitGroup

	// Start multiple worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go doubler(jobs, results, control, &wg)
	}

	// Start job processing timer
	start := time.Now()

	// Send many jobs to the jobs channel
	for i := 0; i < numJobs; i++ {
		jobs <- Job{i, 0}
	}
	close(jobs) // Close the jobs channel as no more jobs will be sent

	// Collect results
	go func() {
		for i := 0; i < numJobs; i++ {
			result := <-results
			fmt.Printf("Processed Job: %d, Result: %d\n", result.data, result.result)
		}
		close(results) // Close results when all jobs are processed
	}()

	// Wait for all workers to finish
	wg.Wait()

	// Signal all workers to exit
	for i := 0; i < numWorkers; i++ {
		control <- DoExit
	}

	// Wait for all workers to acknowledge exit
	for i := 0; i < numWorkers; i++ {
		<-control
	}

	// End job processing timer
	elapsed := time.Since(start)
	fmt.Printf("All jobs processed in %s\n", elapsed)
	fmt.Println("Program exit")
}
