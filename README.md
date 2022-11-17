# monte-carlo-assignment
Sepandar Sepehr

## Prerequisite
Make sure you have the latest version of Go installed

## Running the App
First run `go build observability_app.go` to make sure dependencies are installed

Then run `go run observability_app.go`

Go to `http://localhost:8080/observability` on your browser.

## Implementation Details
`observability_app` is the main runner. It does the following:
* Sets up logger
* Creates context with 30s timeout for the APIs
* Sets up DB. I am using sqlite for local storage for this demo. 
* Sets up ingestion part and adds it to a cron schedule that runs every minute
* Creates handlers for APIs
 