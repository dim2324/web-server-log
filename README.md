# web-server-log
Multithreaded web server log handler

- Reads entries from a CSV file with logs.
- Processes them in parallel using a worker pool (minimum 3 goroutines).
- Filters entries by status code (e.g., only 4xx and 5xx errors).
- Calculates statistics: total number of requests, number of errors, average response time, top IP addresses.
- Outputs results to the console.


Expected project structure:
log-processor/
├── main.go # program entry point
├── processor.go # log processing logic (readLogs, processLogs, filterLogs, calculateStats)
├── testdata/
│ └── logs.csv # test log file (minimum 10-15 entries)
├── go.mod # Go module file
├── README.md # launch instructions with command examples
