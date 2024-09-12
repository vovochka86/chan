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

func doubler(jobs, results chan Job, wg *sync.WaitGroup) {
	defer wg.Done() // Notify when goroutine finishes

	for job := range jobs { // Continuously read from the jobs channel
		results <- Job{data: job.data, result: job.data * 2}
	}
}

func main() {
	numJobs := 5000            // Number of jobs for stress testing
	numWorkers := 10           // Number of concurrent workers
	jobs := make(chan Job, 100) // Buffered channel for job handling
	results := make(chan Job, 100)

	var workerWg sync.WaitGroup // Synchronize worker goroutines

	// Start multiple worker goroutines
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go doubler(jobs, results, &workerWg)
	}

	// Start job processing timer
	start := time.Now()

	// Send many jobs to the jobs channel
	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- Job{i, 0}
		}
		close(jobs) // Close the jobs channel as no more jobs will be sent
	}()

	// Collect and process results concurrently
	var resultWg sync.WaitGroup
	resultWg.Add(1)
	go func() {
		defer resultWg.Done() // Notify when results processing is done
		for i := 0; i < numJobs; i++ {
			result := <-results
			fmt.Printf("Processed Job: %d, Result: %d\n", result.data, result.result)
		}
	}()

	// Wait for all workers to finish processing jobs
	workerWg.Wait()

	// Close the results channel after all workers are done
	close(results)

	// Wait for the results processing to finish
	resultWg.Wait()

	// End job processing timer
	elapsed := time.Since(start)
	fmt.Printf("All jobs processed in %s\n", elapsed)
	fmt.Println("Program exit")
}
