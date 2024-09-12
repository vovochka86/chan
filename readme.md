Key Enhancements for Stress Testing:
Increased jobCount: The number of jobs has been increased to 1000.
Multiple Goroutines: The number of worker goroutines (workerCount) is set to 5 to test concurrency.
Random Processing Time: Each job is simulated with a random processing time of up to 500 milliseconds to create varying workload intensities.
Dynamic Timeout: The timeout for the program is set based on the number of jobs, which ensures that the system can handle larger workloads but still enforces a reasonable timeout.
Performance Monitoring: The program tracks and prints the total time taken to process all jobs.
