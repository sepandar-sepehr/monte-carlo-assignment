# monte-carlo-assignment
Author: Sepandar Sepehr

## Prerequisite
Make sure you have the latest version of Go installed

## Running the App
First run `go build observability_app.go` to make sure dependencies are installed

Then run `go run observability_app.go`

Go to `http://localhost:8080/observability` on your browser.

You should see a dropdown, where you can select from the list of available quotes. 
It will show you last 24h data as well as their rank in that group.
You can see the sample screenshot in the web folder
![screenshot](https://raw.githubusercontent.com/sepsep/monte-carlo-assignment/main/web/Web%20Screenshot.png)

This will make an API call to `http://localhost:8080/market/24h_price` with query params like: `?exchange=coinbase-pro&from=eth&to=eur`
The response includes all the fetched data points for the last 24 hours for that quote as well as that quote's rank (displayed on the top of graph).

## Implementation Details
`observability_app` is the main runner. It does the following:
* Sets up logger.
* Creates context with 30s timeout for the APIs.
* Sets up DB. I am using sqlite for local storage for this demo. 
* Sets up ingestion part and adds it to a cron schedule that runs every minute.
  * We get this data for a list of `SupportedQuotes` from external client and stores them in the DB (using Repository pattern)
* Sets up rank calculator and adds it to another cron schedule that also runs every minute.
  * This calculator uses the data stored in DB from the above job to calculate the ranks 
* Creates handlers for APIs. There are two APIs here
  * One for serving the front-end page
  * One for getting last 24 hours for a quote that is mentioned above

# Next Steps
## Scaling
### Metrics
Right now I am supporting only one exchange (`coinbase-pro`), which is hardcoded. 
This needs to be generalized properly, so we can support more exchange sources based on different quotes.

Moreover, I support only 3 quotes but this can extend to many more by adding them to `SupportedQuotes`. 
However, we need to group them properly, so we rank different quotes based on their grouping.

Another limitation in expansion right now is how my front-end works. The list of supported quotes is hardcoded right now.
I need to add another API to fetch supported quotes and use that for displaying. 
This has to be grouped properly, so they can be shown properly and their ranks make sense.

### Increased Load
If we have multiple users requesting data, we mainly need to make sure our server can handle the traffic properly.
Using a Load Balancer with more instances can help. Or ideally moving the handler part to a serverless architecture such as 
relying on AWS Lambda helps with scaling it without dedicating instances to it. 
We just need to make sure our DB can support the extra load. Moving our storage to a NoSQL storage such as 
DynamoDB can help with that since we are not using any joins or require transactions but we need to
make sure latency is not spiked.

For increasing sampling rate, it is better to change the API based pulling to a push based model.
Cryptowat that is used here already provides WebSocket that we can use to ingest data using a push based method.
We can set up a service for ingestion separately and use it to subscribe to all the trades, then publish them to a 
queue. E.g., we can use a Kafka stream and have a Flink service running to use trade data to calculate
quotes to be stored in the DB. 

## Testing
I completely skipped testing for now. Of course we need to add unit test for all the files in 
this code base (perhaps just skip the models). 

Other than that, we need integration tests if we were going to productionize it and add run them after 
changes to our development environment before pushing it to staging or production. Integration tests 
should test our APIs to provide reasonable data given the input.

## Threshold alerting Feature
A quick solution to this is building everything on top of what we have right now and query the DB
every time we ingest a new value for a metric. We then send an alert if the new metric crosses the
newly calculated threshold. 

To reduce the load on the DB, we can make some improvements. One way is to use a cache for recent
data points, so we get them from cache instead of DB. Another solution is, like mentioned above in 
dealing with increased load, to use a Kafka stream and Flink jobs for that since we can properly 
deal with the moving window using Flink and Kafka and implement the calculations/functions in Flink.
The results will be new values that can be published back into another Kafka topic to be ingested 
and stored in DB as well as being used by a new alerting service to notify subscribers of that monitor. 
