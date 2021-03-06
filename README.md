# go-API-basic

A RESTful API template (built with Go)

The goal of this project is to make an example/template of a relational database-backed REST HTTP API that has characteristics needed to ensure success in a high volume environment. I'm gearing this towards beginners, as I struggled with a lot of this over the past couple of years and would like to help others getting started.

## API Walkthrough

The following is an in-depth walkthrough of this project. This walkthrough has a stupid amount of detail. This is a demo API, so the "business" intent of it is to support basic CRUD (**C**reate, **R**ead, **U**pdate, **D**elete) operations for a movie database.

## Installation

TL;DR - just show me how to install and run the code.

```bash
git clone https://github.com/gilcrest/go-api-basic.git
```

There is a mock implementation which you can run without setting up or having to connect to a database. To validate your installation and run the mock installation, do the following:

Build the code from the root directory

```bash
go build -o server
```

> This sends the output of `go build` to a file called `server` in the same directory.

Execute the file

```bash
./server -env=local -datastore=mock
```

You should see something similar to the following:

```bash
$ ./server -env=local -datastore=mock
{"time":1577934052,"message":"Logging Level set to error"}
{"time":1577934052,"message":"Running, connected to the Local environment, datastore is set to Mock"}
```

You can then execute all four operations (create, read, update, delete) using the following cURL commands.

### cURL Commands to Call API

**Create** - use the POST HTTP verb at `/api/v1/movies`:

```bash
curl --location --request POST 'http://127.0.0.1:8080/api/v1/movies' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Repo Man",
    "year": 1984,
    "rated": "R",
    "release_date": "Mar 02 1984",
    "run_time": 92,
    "director": "Alex Cox",
    "writer": "Alex Cox"
}
'
```

**Read (All Records)** - use the GET HTTP verb at `/api/v1/movies`:

```bash
curl --location --request GET 'http://127.0.0.1:8080/api/v1/movies' \
--data-raw ''
```

**Read (Single Record)** - use the GET HTTP verb and the movie "external ID" as the unique identifier. I try to never expose primary keys, so I use something like an external id. When calling a mock service, this can be anything and it will be echoed in the response e.g. `/api/v1/movies/:extl_id`

```bash
curl --location --request GET 'http://127.0.0.1:8080/api/v1/movies/SKuGy0k6VAojqes40Na' \
--data-raw ''
```

**Update** - use the PUT HTTP verb at `/api/v1/movies/:extl_id`

```bash
curl --location --request PUT 'http://127.0.0.1:8080/api/v1/movies/SKuGy0k6VAojqes40Nga' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Repo Them/They",
    "year": 1984,
    "rated": "R",
    "release_date": "Mar 02 1984",
    "run_time": 92,
    "director": "Alex Cox",
    "writer": "Alex Cox"
}
'
```

**Delete** - use the DELETE HTTP verb at `/api/v1/movies/:extl_id`

### Database Setup

For a non-mocked, persisted implementation using PostgreSQL, you need to setup the database. PostgreSQL needs to be installed locally or one can connect to [Cloud SQL for PostgreSQL](https://cloud.google.com/sql/docs/postgres/) using the [Google Cloud SQL Proxy](https://cloud.google.com/sql/docs/postgres/sql-proxy) or [Cloud Run](https://cloud.google.com/run/) - or any other place you can run PostgreSQL, but instructions are provided only for these three as part of this repo. In all cases, the database will need to be setup and the [demo.ddl](https://github.com/gilcrest/go-api-basic/blob/master/demo.ddl) script found at the root of this project must be executed successfully. You can change the db and schema name, but you'll need to change those in the code as well when doing so. The script creates a simple table and function for inserts.

In addition, environment variables need to be setup for the database of your choice, examples are below (you can setup all three options and choose at build/run time based on a flag).

#### Database Connection Environment Variables

##### Local PostgreSQL

```bash
#Local Postgres DB Name
export PG_APP_DBNAME="fakeDBName"
#Local Postgres DB Username
export PG_APP_USERNAME="fakeDBUsername"
#Local Postgres DB Password
export PG_APP_PASSWORD="fakeDBPassword"
#Local Postgres DB Host
export PG_APP_HOST="localhost"
#Local Postgres DB Port
export PG_APP_PORT="5432"
```

##### Google Cloud SQL Proxy

```bash
#GCP Postgres Cloud Proxy DB Name
export PG_GCP_CP_DBNAME="fakeDBName"
#GCP Postgres Cloud Proxy DB Username
export PG_GCP_CP_USERNAME="postgres"
#GCP Postgres Cloud Proxy DB Password
export PG_GCP_CP_PASSWORD="fakeDBPassword"
#GCP Postgres Cloud Proxy DB Host
export PG_GCP_CP_HOST="localhost"
#GCP Postgres Cloud Proxy DB Port
export PG_GCP_CP_PORT="3307"
```

##### Google Cloud Run

```bash
# GCP Postgres DB Host - for GCP Cloud SQL, use given "instance
# connection name" (which is a concatenation of
# project-id:region:db-instance-id, something like
# "go-api-basic-project:us-east1:go-api-basic-db)
export PG_GCP_HOST="/cloudsql/go-api-basic-project:us-east1:go-api-basic-db"
#GCP Postgres DB Port
export PG_GCP_PORT="5432"
#GCP Postgres DB Name
export PG_GCP_DBNAME="fakeDBName"
#GCP Postgres DB Username
export PG_GCP_USERNAME="postgres"
#GCP Postgres DB Password
export PG_GCP_PASSWORD="fakeDBPassword"
```

### Build and Run

After completing database and environment variable setup, you need to use different build and run commands depending on the environment/datastore. The `-env` flag in the commands below is largely just informational and does not impact the code in any way - the idea is that you can have local, qa, staging, production or whatever values you want and you may want your code to behave in certain ways, but for now it's just symbolic. The `-datastore` flag is the real driver of what is being used for persistence (or as above for mocked data).

#### Local PostgreSQL Datastore - Build and Run

Build the code from the root directory:

```bash
go build -o server
```

> This sends the output of `go build` to a file called `server` in the same directory.

Execute the file from the build step:

```bash
./server -env=local -datastore=local
```

You should see output that looks like the following:

```bash
 ./server -env=local -datastore=local
{"time":1579191386,"message":"Logging Level set to error"}
{"time":1579191386,"message":"Running, connected to the Local environment, datastore is set to Local"}
```

You can then run the same cURL commands as above and validate the changes in your local database.

#### Google Cloud SQL for PostgreSQL Proxy - Build and Run

To run your code locally against the Google Cloud SQL Proxy, there are several dependencies.

1. Setup a project on [Google Cloud](https://cloud.google.com/)
1. Setup the `gcloud` command-line tools locally [link](https://cloud.google.com/sdk/gcloud/)
1. Setup a `PostgreSQL Cloud SQL Database` in that project [link](https://cloud.google.com/sql/docs/postgres/)
1. Setup Google Cloud SQL Proxy [link](https://cloud.google.com/sql/docs/postgres/sql-proxy)

Assuming those dependencies are in place and your database is running, you can start your Cloud SQL Proxy by navigating to the directory you installed the `cloud_sql_proxy` tool and using the following:

```bash
./cloud_sql_proxy -instances=go-api-basic-project:us-east1:go-api-basic-db=tcp:3307
```

This tells the `cloud_sql_proxy` tool to connect to the instance above and listen locally on port 3307. You should see something like the following:

```bash
$ ./cloud_sql_proxy -instances=go-api-basic-project:us-east1:go-api-basic-db=tcp:3307
2020/01/16 10:59:26 current FDs rlimit set to 10240, wanted limit is 8500. Nothing to do here.
2020/01/16 10:59:28 Listening on 127.0.0.1:3307 for go-api-basic-project:us-east1:go-api-basic-db
2020/01/16 10:59:28 Ready for new connections
```

If you haven't already done so, build the code from the project's root directory:

```bash
go build -o server
```

> This sends the output of `go build` to a file called `server` in the same directory.

Execute the file from the build step:

```bash
./server -env=local -datastore=gcpcp
```

You should see output that looks like the following:

```bash
$ ./server -env=local -datastore=gcpcp
{"time":1579191760,"message":"Logging Level set to error"}
{"time":1579191760,"message":"Running, connected to the Local environment, datastore is set to Google Cloud SQL through the Google Cloud Proxy"}
```

You can then run the same cURL commands as above and validate the changes in your Google Cloud SQL PostgreSQL database.

#### Google Cloud Run - Build and Run

In order to run this project on `Google Cloud Run`, there are several dependencies:

1. Setup a project on [Google Cloud](https://cloud.google.com/)
1. Setup the `gcloud` command-line tools locally [link](https://cloud.google.com/sdk/gcloud/)
1. Setup a `PostgreSQL Cloud SQL Database` in that project [link](https://cloud.google.com/sql/docs/postgres/)
1. You should also research Cloud Run to understand what it is [link](https://cloud.google.com/run/)

In order to deploy this code as a container, it needs to be built to a container registry. I am using the [Google Container Registry](https://cloud.google.com/container-registry/docs/) aka GCR to store my built container images. You can use the `gcloud` command line tool to build to this registry using the following:

```bash
gcloud builds submit --tag gcr.io/go-api-basic/demo
```

You should see a lot of output as part of this, however, the end result should say SUCCESS and point to the image you just uploaded.

In order to deploy the container to Cloud Run, you need to have the Cloud Run environment variables from above setup with your configurations an run the following command:

```bash
gcloud run deploy movies --image gcr.io/go-api-basic/demo --platform managed \
--add-cloudsql-instances go-api-basic-project:us-east1:go-api-basic-db \
--set-env-vars \
INSTANCE-CONNECTION-NAME=go-api-basic-project:us-east1:go-api-basic-db,\
PG_GCP_DBNAME=${PG_GCP_DBNAME},\
PG_GCP_USERNAME=${PG_GCP_USERNAME},\
PG_GCP_PASSWORD=${PG_GCP_PASSWORD},\
PG_GCP_HOST=${PG_GCP_HOST},\
PG_GCP_PORT=${PG_GCP_PORT}
```

This command takes the container we just put to the Google Container Registry and deploys it to Cloud Run. It also registers which database to use with the `--add-cloudsql-instances` flag. After deploying, you should see something like the following:

```bash
Deploying container to Cloud Run service [movies] in project [go-api-basic-project] region [us-east1]
✓ Deploying... Done.
  ✓ Creating Revision...
  ✓ Routing traffic...
Done.
Service [movies] revision [movies-00001-abc] has been deployed and is serving 100 percent of traffic at https://movies-2fake5enya-ue.a.run.app
```

You can then execute all the same cURL commands as above, you just need to change the base URL to the one given after the deployment. As an example, the POST cURL becomes:

```bash
curl --location --request POST 'https://movies-2fake5enya-ue.a.run.app/api/v1/movies' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Repo Man",
    "year": 1984,
    "rated": "R",
    "release_date": "Mar 02 1984",
    "run_time": 92,
    "director": "Alex Cox",
    "writer": "Alex Cox"
}
'
```

That's it - Google Cloud Run is really cool - we now have a service running with a proper SSL cert and everything!

## Project Walkthrough

### Errors

Before even getting into the full walkthrough, I wanted to review the [errs module](https://github.com/gilcrest/errs) and the approach taken for error handling. The `errs module` is basically a carve out of the error handling used in the [upspin library](https://github.com/upspin/upspin/tree/master/errors) with some tweaks and additions I made for my own needs. Rob Pike has a [fantastic post](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html) about errors and the Upspin implementation. I've taken that and added my own twist.

The general idea for error handling throughout this example API is to always raise an error using the `errs.E` function as seen in this simple error handle below. `errs.E` is neat - you can pass in any one of a number of approved types and the function helps form the error. In all error cases, pass `errs.Op` (for the **Op**eration) as the `errs.E` function helps build a pseudo stack trace for the error as it goes up through the code. Here's a snippet showing a typical, simple example of using the errs.E function.

```go
func NewServer(name env.Name, lvl zerolog.Level) (*Server, error) {
    const op errs.Op = "server/NewServer"

    // call constructor for Env struct from env module
    env, err := env.NewEnv(name, lvl)
    if err != nil {
        return nil, errs.E(op, err)
    }
```

The following snippet shows a more robust validation example. In it, you'll notice that if you need to define your own error quickly, you can just use a string and that becomes the error string as well.

```go
func (m *Movie) validate() error {
    const op errs.Op = "movie/Movie.validate"

    switch {
    case m.Title == "":
        return errs.E(op, errs.Validation, errs.Parameter("Title"), errs.MissingField("Title"))
    case m.Year < 1878:
        return errs.E(op, errs.Validation, errs.Parameter("Year"), "The first film was in 1878, Year must be >= 1878")
```

In the above snippet, the errs.MissingField function used to validate missing input on fields comes from [this Mat Ryer post](https://medium.com/@matryer/patterns-for-decoding-and-validating-input-in-go-data-apis-152291ac7372) and is pretty handy.

```go
// MissingField is an error type that can be used when
// validating input fields that do not have a value, but should
type MissingField string

func (e MissingField) Error() string {
    return string(e) + " is required"
}
```

As stated before, as errors go up the stack from whatever depth of code they're in, the operation is captured and added to the error string as a pseudo stack trace that is super helpful for debugging. However, I don't want this type of internal stack information exposed to end users in the response - I only want the error message. As such, just prior to shipping the response, I log the error (to capture the stack info) and call a custom function I built called `errs.RE` (**R**esponse **E**rror). This function effectively strips the stack information and just sends the original error message along with whatever http status code you select as well as whatever errs.Kind, Code or Parameter you choose to set. The `RE` function returns an error of type `errs.HTTPErr`.  An example of error handling at the highest level (from the POST handler) is below:

```go
// Call the create method of the Movie object to validate and insert the data
err = movie.Create(ctx, s.Logger, tx)
if err != nil {
    // log error
    s.Logger.Error().Err(err).Msg("")
    // Type assertion is used - all errors should be an *errs.Error type
    // Use Kind, Param, Code and Error from lower level errors to populate RE (Response Error)
    if e, ok := err.(*errs.Error); ok {
        err := errs.RE(http.StatusBadRequest, e.Kind, e.Param, e.Code, err)
        errs.HTTPError(w, err)
        return
    }

    // if falls through type assertion, then serve an unanticipated error
    err := errs.RE(http.StatusInternalServerError, errs.Unanticipated)
    errs.HTTPError(w, err)
    return
}
```

The final statement above before returning the errors is a call to the `errs.HTTPError` function. This function determines if an error is of type `errs.HTTPErr` and if so, forms the error json - the response body will look something like this:

```json
{
    "error": {
        "kind": "input_validation_error",
        "param": "Year",
        "message": "The first film was in 1878, Year must be >= 1878"
    }
}
```

## EVERYTHING BELOW HERE IS OUTDATED AND NEEDS TO BE UPDATED - I'M WORKING ON IT ASAP (2/4/2020)

### Main API Module

The [Main API/server layer module](https://github.com/gilcrest/go-api-basic) is the starting point and as such has the main package/function within the cmd directory. In it, I'm checking for both log level and environment command line flags and running them through 2 simple functions (`logLevel` and `envName`) to get the correct string value given the flag input. I'm using [zerolog](https://github.com/rs/zerolog) throughout my modules as the logger.

```go
func main() {

    // loglvl flag allows for setting logging level, e.g. to run the server
    // with level set to debug, it'd be: ./server loglvl=debug
    loglvlFlag := flag.String("loglvl", "error", "sets log level")

    // env flag allows for setting environment, e.g. Production, QA, etc.
    // example: env=dev, env=qa, env=stg, env=prod
    envFlag := flag.String("env", "dev", "sets log level")

    flag.Parse()

    // get appropriate zerolog.Level based on flag
    lvl := logLevel(loglvlFlag)
    log.Log().Msgf("Logging Level set to %s", lvl)

    // get appropriate env.Name based on flag
    eName := envName(envFlag)
    log.Log().Msgf("Environment set to %s", eName)
```

Next in the `main`, I'm calling the `NewServer` function from the `server` package to construct a new server using the environment name and logging level.

```go
    // call constructor for Server struct
    server, err := server.NewServer(eName, lvl)
    if err != nil {
        log.Fatal().Err(err).Msg("")
    }
```

The server's multiplex router is Gorilla, which is registered as the handler for http.Handle.

```go
    // handle all requests with the Gorilla router
    http.Handle("/", server.Router)
```

Finally, http.ListenAndServe is run to truly start the server and listen for incoming requests.

```go
    // ListenAndServe on port 8080, not specifying a particular IP address
    // for this particular implementation
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal().Err(err).Msg("")
    }
```

Let's dig into the `server.NewServer` constructor function from above. In the `server` package, there is also a `Server` struct, which is composed of a pointer to the `Env` struct from the [env module](https://github.com/gilcrest/env), as below:

```go
// Server struct contains the environment (env.Env) and additional methods
// for running our HTTP server
type Server struct {
    *env.Env
}
```

The `Server` struct uses [type embedding](https://golang.org/doc/effective_go.html#embedding), which allows it to take on all the properties of the `Env` struct from the `env module` (things like database setup, logger, multiplexer, etc.), but also allows for extending the struct with API-specific methods for routing logic and handlers.

The `Env` struct in the env module has the following structure:

```go
// Env struct stores common environment related items
type Env struct {
    // Environment Name (e.g. Production, QA, etc.)
    Name Name
    // multiplex router
    Router *mux.Router
    // Datastore struct containing AppDB (PostgreSQL),
    // LogDb (PostgreSQL) and CacheDB (Redis)
    DS *datastore.Datastore
    // Logger
    Logger zerolog.Logger
}
```

Within the `server.NewServer` constructor function, I'm calling the `NewEnv` function of the `env module`. The `env.NewEnv` constructor function does all the database setup, starts the logger as well as initializes the gorilla multiplexer. Fore more information on what's happening inside the `env.NewEnv` function, reference the [env module Readme](https://github.com/gilcrest/env)

```go
// NewServer is a constructor for the Server struct
// Sets up the struct and registers routes
func NewServer(name env.Name, lvl zerolog.Level) (*Server, error) {
    const op errs.Op = "server/NewServer"

    // call constructor for Env struct from env module
    env, err := env.NewEnv(name, lvl)
    if err != nil {
        return nil, errs.E(op, err)
    }
```

After getting `env.Env` back, I embed it in the Server struct

```go
    // Use type embedding to make env.Env struct part of Server struct
    server := &Server{env}
```

Finally, I call the `server.routes` method and return the server.

```go
    // routes registers handlers to the Server router
    err = server.routes()
    if err != nil {
        return nil, errs.E(op, err)
    }

    return server, nil

```

Inside the `server.routes` method, first I pull out the app database from the server struct to pass into my `servertoken` handler.

```go
// routes registers handlers to the router
func (s *Server) routes() error {
    const op errs.Op = "server/Server.routes"

    // Get App Database for token authentication
    appdb, err := s.DS.DB(datastore.AppDB)
    if err != nil {
        return errs.E(op, err)
    }
```

Next, the URL path and handlers are register to the router embedded in the server.

```go
    s.Router.Handle("/v1/movie",
        alice.New(
            s.handleStdResponseHeader,
            servertoken.Handler(s.Logger, appdb)).
            ThenFunc(s.handlePost())).
        Methods("POST").
        Headers("Content-Type", "application/json")
```

The `Methods("POST").` means this route will only take POST request, and for REST this means we're looking at our Create method of the (CRUD) we talked about above. Other methods (Read(GET), Update(PUT), and Delete(DELETE)) will be documented later. The `Headers("Content-Type", "application/json")` means that this route requires that this request header be present.

To go through the `v1/movie` Handle registration item by item - [Justinas Stankevičius' alice library](https://github.com/justinas/alice) is being used to make middleware chaining easier.

Next, the first middleware in the chain above `s.handleStdResponseHeader` simply adds standard response headers. As of now, it's just the `Content-Type:application/json` header, but it's an easy place to other headers one may deem standard.

```go
// handleStdResponseHeader middleware is used to add standard HTTP response headers
func (s *Server) handleStdResponseHeader(h http.Handler) http.Handler {
    return http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            w.Header().Add("Content-Type", "application/json")
            h.ServeHTTP(w, r) // call original
        })
}
```

Next, the second middleware in the chain above (`servertoken.Handler`) validates that the caller of the API is authorized. Details for what's happening in this `servertoken module` can be found [here](https://github.com/gilcrest/servertoken).

The final handler, `s.handlePost`, which hangs off the `Server` struct mentioned above, is meant to handle `POST` requests to the route registered for the API ("/v1/movie").

Inside the `handlePost` method, a private request and response struct are defined. I like this technique as it gives me complete control over what is coming and going from the API.

```go
// handlePost handles POST requests for the /movie endpoint
// and creates a movie in the database
func (s *Server) handlePost() http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {

        // request is the expected service request fields
        type request struct {
            Title    string `json:"Title"`
            Year     int    `json:"Year"`
            Rated    string `json:"Rated"`
            Released string `json:"ReleaseDate"`
            RunTime  int    `json:"RunTime"`
            Director string `json:"Director"`
            Writer   string `json:"Writer"`
        }

        // response is the expected service response fields
        type response struct {
            request
            CreateTimestamp string `json:"CreateTimestamp"`
        }
```

The request is Decoded into an instance of the request struct.

```go
        // Declare rqst as an instance of request
        // Decode JSON HTTP request body into a Decoder type
        // and unmarshal that into rqst
        rqst := new(request)
        err := json.NewDecoder(req.Body).Decode(&rqst)
        defer req.Body.Close()
        if err != nil {
            err = errs.RE(http.StatusBadRequest, errs.InvalidRequest, err)
            errs.HTTPError(w, err)
            return
        }
```

The request is mapped to the business struct (`movie.Movie`) from the [movie module](https://github.com/gilcrest/movie).

```go
// Movie holds details of a movie
type Movie struct {
    Title    string
    Year     int
    Rated    string
    Released time.Time
    RunTime  int
    Director string
    Writer   string
    dbaudit.Audit
}
```

As part of the mapping, quick input validations around date formatting are done (see time.Parse below)

```go
        // dateFormat is the expected date format for any date fields
        // in the request
        const dateFormat string = "Jan 02 2006"

        // declare a new instance of movie.Movie
        movie := new(movie.Movie)
        movie.Title = rqst.Title
        movie.Year = rqst.Year
        movie.Rated = rqst.Rated
        t, err := time.Parse(dateFormat, rqst.Released)
        if err != nil {
            err = errs.RE(http.StatusBadRequest,
                errs.Validation,
                errs.Code("invalid_date_format"),
                errs.Parameter("ReleaseDate"),
                err)
            errs.HTTPError(w, err)
            return
        }
        movie.Released = t
        movie.RunTime = rqst.RunTime
        movie.Director = rqst.Director
        movie.Writer = rqst.Writer
```

The context is pulled from the incoming request and a database transaction is started using the AppDB from the Server struct

```go
        // retrieve the context from the http.Request
        ctx := req.Context()

        // get a new DB Tx from the PostgreSQL datastore within the server struct
        tx, err := s.DS.BeginTx(ctx, nil, datastore.AppDB)
        if err != nil {
            err = errs.RE(http.StatusInternalServerError, errs.Database)
            errs.HTTPError(w, err)
            return
        }
```

The `Create` method of the `movie.Movie` struct is called using the context and database transaction from above as well as the Logger from the server. The error handling is important here, but it is discussed at length in the [errors section](#Errors). For more information on what's happening inside the `movie module` check the Readme [here](https://github.com/gilcrest/movie). In summary though, the `movie module` has the "business logic" for the API - it is doing deeper input validations, exercising any business rules and creating/reading/updating/deleting the data in the database (as well as handling commit or rollback).

```go
        // Call the create method of the Movie object to validate and insert the data
        err = movie.Create(ctx, s.Logger, tx)
        if err != nil {
            // log error
            s.Logger.Error().Err(err).Msg("")
            // Type assertion is used - all errors should be an *errs.Error type
            // Use Kind, Param, Code and Error from lower level errors to populate RE (Response Error)
            if e, ok := err.(*errs.Error); ok {
                err := errs.RE(http.StatusBadRequest, e.Kind, e.Param, e.Code, err)
                errs.HTTPError(w, err)
                return
            }

            // if falls through type assertion, then serve an unanticipated error
            err := errs.RE(http.StatusInternalServerError, errs.Unanticipated)
            errs.HTTPError(w, err)
            return
        }
```

If we got this far, the db transaction has been created/committed - we can consider this transaction successful and return a response. An instance of the response struct is initialized and populated with data from the `movie.Movie` struct and the response is encoded and sent back to the caller!

```go
        // create a new response struct and set Audit and other
        // relevant elements
        resp := new(response)
        resp.Title = movie.Title
        resp.Year = movie.Year
        resp.Rated = movie.Rated
        resp.Released = movie.Released.Format(dateFormat)
        resp.RunTime = movie.RunTime
        resp.Director = movie.Director
        resp.Writer = movie.Writer
        resp.CreateTimestamp = movie.CreateTimestamp.Format(time.RFC3339)

        // Encode response struct to JSON for the response body
        json.NewEncoder(w).Encode(*resp)
        if err != nil {
            err = errs.RE(http.StatusInternalServerError, errs.Internal)
            errs.HTTPError(w, err)
            return
        }
```

## Thanks / Attribution

I should say that most of the ideas I'm presenting here are not my own - I learned them from reading a number of books and blogs from extremely talented individuals. Here is a list (in no particular order) of influences:

- [Rob Pike](https://twitter.com/rob_pike?lang=en)
- [Jaana B Dogan](https://twitter.com/rakyll?ref_src=twsrc%5Egoogle%7Ctwcamp%5Eserp%7Ctwgr%5Eauthor)
- [Mat Ryer](https://medium.com/@matryer)
- [Jon Calhoun](https://www.calhoun.io/about)
- [Matt Silverlock](https://twitter.com/elithrar?lang=en)
- [Alex Edwards](https://www.alexedwards.net/)

Questions/Concerns? Want more detail? Feel free to [open an issue](https://github.com/gilcrest/go-api-basic/issues) and label it appropriately. Thanks!
