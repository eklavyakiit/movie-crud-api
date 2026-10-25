
# 🎬 Movie REST API in Go

A beginner-friendly REST API built using **Golang, `net/http`, JSON, and Gorilla Mux**.

This project is my first step towards understanding backend engineering with Go. Instead of just learning the syntax of the language, the idea here is to understand how Go can actually be used to build a server, receive HTTP requests, process data, and send responses back to clients.

The project manages a collection of movies. Users can fetch all movies, retrieve a particular movie, add new movies, edit existing movies, and delete movies using HTTP requests.

The goal isn't to build something production-ready right away. The goal is to understand what is happening behind the scenes.

---

## 1. What are we actually building?

Imagine opening a movie application. You can see a list of movies, click on one to see its details, add a new movie, edit its information, or remove it.

Usually, a frontend application communicates with a backend server to perform these operations. Our project focuses on building that backend.

For now, we don't have a frontend or a database. We have a Go server that exposes a set of URLs called **API endpoints**. A client such as Postman can send requests to these endpoints and receive responses.

For example:

- `GET /movies` → Give me all movies.
- `GET /movies/3` → Give me the movie whose ID is 3.
- `POST /movies` → Add a new movie.
- `PUT /movies/3` → Replace the data of movie 3.
- `DELETE /movies/3` → Delete movie 3.

These operations are commonly grouped under the term **CRUD**:

| Letter | Operation | HTTP method |
|---|---|---|
| C | Create | POST |
| R | Read | GET |
| U | Update | PUT |
| D | Delete | DELETE |

That's the core of this project: building a small CRUD API using Go.

## 2. Before diving into the code: what is an API?

API stands for **Application Programming Interface**. Think of it as a defined way for two pieces of software to communicate.

Suppose a frontend wants to display every movie stored by our server. It doesn't need to know how our Go code stores the movies internally. It simply sends a request to the appropriate endpoint.

Our server receives that request, executes the relevant function, and returns the result.

The communication looks like this:

```text
Client (Postman / Frontend)
          |
          | HTTP Request
          v
     Go HTTP Server
          |
          v
      API Handler
          |
          v
     Movie Slice
          |
          v
      JSON Response
          |
          v
         Client
```

The client and server communicate using HTTP, and JSON is the format we use to represent movie data.

### What is REST?

REST stands for Representational State Transfer. It is an architectural style commonly used to design web APIs.

In our project, movies are resources. The URL identifies the resource, and the HTTP method describes what we want to do with it.

For example, `/movies/3` identifies movie 3. A `GET` request reads it, a `PUT` request updates it, and a `DELETE` request removes it.

This combination of resource-oriented URLs and HTTP methods gives us a simple REST-style API.

## 3. Project structure and imports

Our project can be kept very simple:

```text
movie-api/
│
├── go.mod
├── go.sum
└── main.go
```

- `main.go` contains our application code.
- `go.mod` defines the Go module and its dependencies.
- `go.sum` contains checksums used to verify downloaded dependencies.

The application starts from the `main()` function inside `main.go`.

### Understanding the imports

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
)
```

Let's understand every import.

**`package main`**

Go programs are organized into packages. Declaring `package main` identifies this as an executable program when it includes a `main()` function.

When we run `go run .`, Go builds and executes the program starting at `main()`.

**`encoding/json`**

This is part of Go's standard library. It allows us to convert Go data structures into JSON and JSON data into Go data structures.

For example, it lets our server read movie information sent by a client and return movie information in a format the client understands.

**`fmt`**

Used for formatted input and output. Here, we use:

```go
fmt.Println("Starting server at port 8080")
```

This prints a message in the terminal.

**`log`**

Used for logging. Our application uses:

```go
log.Fatal(http.ListenAndServe(":8080", r))
```

`log.Fatal` prints an error and terminates the program if the server cannot start or if the server function returns an error.

**`net/http`**

This is Go's built-in HTTP package. It provides the tools needed to build an HTTP server, receive requests, write responses, and define HTTP handlers.

**`strconv`**

This package converts between strings and numeric types.

For example:

```go
id := strconv.Itoa(7)
```

This converts the integer `7` into the string `"7"`.

Our movie IDs are strings, so this is useful when generating IDs from numbers.

**`github.com/gorilla/mux`**

Gorilla Mux is an external HTTP routing library. Go's `net/http` already provides basic routing, but Mux makes it convenient to define routes containing variables and restrict routes to particular HTTP methods.

For example:

```go
r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
```

The `{id}` portion is a variable that Mux extracts from the request URL.

Unlike the other imports, Gorilla Mux is not part of Go's standard library. It must be included as a project dependency.

---

## 4. Structs: how we represent a movie

Before understanding the handlers, we need to understand the data they work with.

In Go, a `struct` groups related fields into one data type. Think of it as a blueprint for an object containing multiple pieces of information.

Our project defines two structs:

```go
type Movie struct {
    ID       string   `json:"id"`
    Isbn     string   `json:"isbn"`
    Title    string   `json:"title"`
    Director Director `json:"director"`
}

type Director struct {
    FirstName string `json:"firstname"`
    LastName  string `json:"lastname"`
}
```

### The `type` keyword

```go
type Movie struct { ... }
```

This declares a new type named `Movie`, whose underlying structure is defined by the fields inside the braces.

Whenever we write:

```go
var movie Movie
```

we create a variable named `movie` whose type is `Movie`.

We can then access its fields using dot notation:

```go
movie.Title = "Interstellar"
movie.Isbn = "100007"
```

### Why do we have a separate `Director` struct?

A movie has a director, and a director has a first name and a last name. Instead of storing the director's information as one long string, we represent it using a separate struct.

```go
type Director struct {
    FirstName string `json:"firstname"`
    LastName  string `json:"lastname"`
}
```

We then use it inside `Movie`:

```go
Director Director
```

The first `Director` is the field name, and the second is the type.

This is called **struct composition**: one struct contains another struct as a field.

For example:

```go
movie := Movie{
    ID:    "1",
    Isbn:  "100001",
    Title: "3 Idiots",
    Director: Director{
        FirstName: "Rajkumar",
        LastName:  "Hirani",
    },
}
```

Notice that the nested `Director` struct has its own set of braces.

This structure makes the data organized and easy to work with.

### What are JSON struct tags?

Look at this field:

```go
Title string `json:"title"`
```

The part inside backticks is a struct tag. Go's JSON encoder and decoder use it to determine the corresponding JSON field name.

The Go field is called `Title`, but the JSON field is called `title`.

Our JSON representation looks like this:

```json
{
    "id": "1",
    "isbn": "100001",
    "title": "3 Idiots",
    "director": {
        "firstname": "Rajkumar",
        "lastname": "Hirani"
    }
}
```

**Important syntax detail:** the correct tag is `` `json:"title"` ``, without a space after `json:`.

If the tag is malformed, the JSON package may ignore the intended tag and fall back to its default field-name behavior.

Also notice that `Movie` contains a `Director` value, not a `*Director` pointer. That means we initialize it using `Director{...}`, rather than `&Director{...}`.

---

## 5. The `movies` slice: where our data lives

```go
var movies []Movie
```

This line is incredibly important because almost every endpoint interacts with this variable.

Let's break it down.

- `var` declares a variable.
- `movies` is the variable name.
- `[]Movie` means a slice whose elements have type `Movie`.

### What is a slice?

A slice is a flexible, dynamically sized view of an underlying array. Unlike a fixed-length array, a slice can grow when we append new elements.

For example:

```go
var movies []Movie

movies = append(movies, Movie{
    ID:    "1",
    Isbn:  "100001",
    Title: "3 Idiots",
    Director: Director{
        FirstName: "Rajkumar",
        LastName:  "Hirani",
    },
})
```

The `append()` function adds the movie to the slice and returns the resulting slice. That's why we assign its result back to `movies`.

We can also add multiple initial movies inside `main()` using the same approach.

### How do we access elements?

Slices use zero-based indexing:

```go
movies[0] // First movie
movies[1] // Second movie
movies[2] // Third movie
```

If we have six movies, their indexes are `0` through `5`.

This matters when we update or delete a movie because we often need its position in the slice.

### The limitation of our storage

Our data is stored in RAM, not in a database.

This means:

- New movies exist only while the program is running.
- Updates modify the in-memory slice.
- Deletions remove elements from the in-memory slice.
- Restarting the server resets the collection to the initial movies defined in `main()`.

This is completely fine for learning CRUD operations. Later, we can replace the slice with a database without fundamentally changing the idea of our API endpoints.

---

## 6. HTTP handlers: the functions that handle requests

An HTTP handler is code that responds to an HTTP request.

Our functions follow this pattern:

```go
func getMovies(w http.ResponseWriter, r *http.Request) {
    // Handle request
}
```

Let's understand the syntax.

**`func`** declares a function.

**`getMovies`** is the function name.

**`w http.ResponseWriter`** is the response writer. We use it to set response headers, write status codes, and send data back to the client.

**`r *http.Request`** is a pointer to the incoming HTTP request. It contains information such as the request method, URL, headers, and body.

Go's `net/http` package calls the handler when a matching request arrives.

### Setting the response content type

Most of our handlers start with:

```go
w.Header().Set("Content-Type", "application/json")
```

This tells the client that the response body contains JSON.

The header is metadata about the response, not the JSON data itself.

We use:

```go
json.NewEncoder(w).Encode(movie)
```

to actually write the movie data into the response.

`Encode()` also writes a trailing newline, which is perfectly normal for JSON responses.

---

## 7. GET: retrieving all movies

Our first endpoint is:

```go
func getMovies(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movies)
}
```

This is probably the simplest handler in the project.

When a client sends:

```http
GET /movies
```

the server runs `getMovies()`.

First, we set the content type to JSON. Then we encode the entire `movies` slice and write it to the response.

Because the slice contains multiple `Movie` structs, the JSON output is an array of movie objects.

For example:

```json
[
    {
        "id": "1",
        "isbn": "100001",
        "title": "3 Idiots",
        "director": {
            "firstname": "Rajkumar",
            "lastname": "Hirani"
        }
    },
    {
        "id": "2",
        "isbn": "100002",
        "title": "Dangal",
        "director": {
            "firstname": "Nitesh",
            "lastname": "Tiwari"
        }
    }
]
```

The real response contains all the movies currently in the slice.

### What is `json.NewEncoder()`?

```go
json.NewEncoder(w).Encode(movies)
```

Think of this as two operations combined:

1. `json.NewEncoder(w)` creates an encoder that writes JSON to `w`.
2. `.Encode(movies)` converts `movies` into JSON and writes it to the response.

The `w` argument connects the encoder to the HTTP response.

---

## 8. GET by ID: retrieving one movie

The second endpoint retrieves a particular movie:

```go
func getMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    params := mux.Vars(r)

    for _, item := range movies {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }

    http.Error(w, "Movie not found", http.StatusNotFound)
}
```

Let's walk through this carefully.

### Step 1: Extract the ID

```go
params := mux.Vars(r)
```

Gorilla Mux examines the route and extracts any path variables.

Our route is:

```go
r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
```

Suppose the client requests `/movies/3`.

Mux gives us a map that conceptually looks like this:

```go
map[string]string{
    "id": "3",
}
```

We retrieve the value using:

```go
params["id"]
```

The result is the string `"3"`.

You don't need to initialize this map yourself. Gorilla Mux creates and populates the route-variable map when it matches the request to the route.

### Step 2: Loop through the slice

```go
for _, item := range movies {
```

The `range` keyword iterates over the elements of a slice.

Here:

- `_` ignores the index because we don't need it.
- `item` represents the current movie.

We could use an index as well, but in this function, we only need to read each movie.

### Step 3: Compare IDs

```go
if item.ID == params["id"] {
```

We check whether the current movie's ID matches the ID from the URL.

If it does, we return that movie:

```go
json.NewEncoder(w).Encode(item)
return
```

The `return` statement stops the function immediately. Without it, the loop could continue and the function might attempt to process further elements.

### Step 4: Handle missing movies

If the loop finishes without finding a match, we execute:

```go
http.Error(w, "Movie not found", http.StatusNotFound)
```

This writes an error response with HTTP status `404`.

That's better than silently returning an empty response because the client can tell that the requested movie doesn't exist.

---

## 9. POST: creating a new movie

Our creation handler receives movie data from the client and appends it to the slice.

```go
func createMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var movie Movie

    if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    maxID := 0

    for _, item := range movies {
        id, err := strconv.Atoi(item.ID)
        if err == nil && id > maxID {
            maxID = id
        }
    }

    movie.ID = strconv.Itoa(maxID + 1)

    movies = append(movies, movie)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(movie)
}
```

This handler introduces several important Go concepts.

### Step 1: Create a variable to hold the incoming movie

```go
var movie Movie
```

We create a `Movie` variable that will be populated using the JSON sent by the client.

At first, its fields have their zero values. For strings, that means empty strings.

### Step 2: Decode the request body

```go
if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
}
```

The client might send this request body:

```json
{
    "isbn": "100007",
    "title": "Interstellar",
    "director": {
        "firstname": "Christopher",
        "lastname": "Nolan"
    }
}
```

`r.Body` contains the incoming JSON data.

`json.NewDecoder(r.Body).Decode(&movie)` reads that data and fills the fields of the `movie` variable.

Notice that we pass `&movie`, a pointer to the variable. The decoder needs to modify that variable.

### What is `if err := ...; err != nil`?

This is a common Go pattern called an *if statement with a short variable declaration*.

```go
if err := someOperation(); err != nil {
    // Handle error
}
```

It does three things:

1. Runs the operation.
2. Stores its error result in `err`.
3. Checks whether the error is not `nil`.

In Go, `nil` represents the absence of a value for certain types, including errors.

If `err != nil`, something went wrong. In our case, the incoming JSON could be malformed, so we return a `400 Bad Request` response.

This pattern is used extensively throughout Go code.

### Step 3: Generate an ID

```go
maxID := 0
```

The `:=` operator declares and initializes a variable. Go infers its type from the assigned value, so `maxID` becomes an integer.

We then examine the IDs of existing movies:

```go
for _, item := range movies {
    id, err := strconv.Atoi(item.ID)

    if err == nil && id > maxID {
        maxID = id
    }
}
```

`strconv.Atoi()` converts a string into an integer. For example, `"6"` becomes `6`.

It returns two values:

- The converted integer.
- An error, if conversion fails.

The condition:

```go
err == nil && id > maxID
```

means that the conversion succeeded and the current ID is greater than the maximum we have seen so far.

The `&&` operator means logical AND: both conditions must be true.

After the loop, we generate the next ID:

```go
movie.ID = strconv.Itoa(maxID + 1)
```

`strconv.Itoa()` converts the integer back to a string.

If the largest ID is `"6"`, the new movie gets `"7"`.

This is a simple learning-friendly strategy. It assumes IDs are numeric and doesn't account for concurrent requests or more complicated ID management. A real application would usually use a database-generated ID or a UUID.

### Step 4: Append the movie

```go
movies = append(movies, movie)
```

This adds the new movie to the slice.

Remember that `append()` returns the resulting slice, which is why we assign it back to `movies`.

### Step 5: Return the new movie

```go
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(movie)
```

`http.StatusCreated` is HTTP status `201`. It indicates that a resource has been created successfully.

We then return the newly created movie as JSON, including the generated ID.

---

## 10. PUT: updating an existing movie

This is the endpoint that lets a user edit a movie.

The client sends a `PUT` request to a URL such as `/movies/3`, with the complete updated movie data in the request body.

Here's the core implementation:

```go
func updateMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    params := mux.Vars(r)
    id := params["id"]

    var updatedMovie Movie

    err := json.NewDecoder(r.Body).Decode(&updatedMovie)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    for index, item := range movies {
        if item.ID == id {
            updatedMovie.ID = item.ID
            movies[index] = updatedMovie

            json.NewEncoder(w).Encode(updatedMovie)
            return
        }
    }

    http.Error(w, "Movie not found", http.StatusNotFound)
}
```

### Step 1: Identify the movie

```go
params := mux.Vars(r)
id := params["id"]
```

The URL tells us which movie should be updated.

For example, `/movies/3` means that the target movie has ID `"3"`.

### Step 2: Decode the updated data

```go
var updatedMovie Movie

err := json.NewDecoder(r.Body).Decode(&updatedMovie)
```

We create a new `Movie` variable and populate it with the data sent by the client.

For example, the client could send:

```json
{
    "isbn": "200003",
    "title": "Shershaah: Director's Cut",
    "director": {
        "firstname": "Vishnuvardhan",
        "lastname": "Rao"
    }
}
```

The new information is stored in `updatedMovie`.

### Step 3: Find the movie and overwrite it

```go
for index, item := range movies {
    if item.ID == id {
        updatedMovie.ID = item.ID
        movies[index] = updatedMovie
        json.NewEncoder(w).Encode(updatedMovie)
        return
    }
}
```

This is the heart of the update operation.

We use `index` because we need to know the position of the existing movie in the slice.

Once we find the matching ID, we preserve it:

```go
updatedMovie.ID = item.ID
```

Why? Because the client should not be able to change the movie's identity just by sending a different ID in the request body. The URL identifies the resource being updated.

Then we replace the existing element:

```go
movies[index] = updatedMovie
```

This overwrites the movie at that position in the slice. We don't have to delete it and append a replacement.

Finally, we encode the updated movie into the response and return from the function.

### PUT versus PATCH

Our implementation uses `PUT`, which normally means replacing the resource with the representation supplied by the client.

That means the client should send the entire movie object, including all the fields it wants to retain.

If a client sends only:

```json
{
    "title": "New Title"
}
```

the other fields in `updatedMovie` will remain at their zero values. The current function will replace the old movie with that incomplete data.

If we wanted to edit only one field while preserving all the others, we would implement partial-update behavior, commonly through `PATCH`.

For now, full replacement is a good way to learn the fundamentals.

---

## 11. DELETE: removing a movie

Our delete handler removes a movie from the slice using its ID.

```go
func deleteMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    params := mux.Vars(r)

    for index, item := range movies {
        if item.ID == params["id"] {
            movies = append(movies[:index], movies[index+1:]...)

            json.NewEncoder(w).Encode(movies)
            return
        }
    }

    http.Error(w, "Movie not found", http.StatusNotFound)
}
```

The important line is:

```go
movies = append(movies[:index], movies[index+1:]...)
```

This can look confusing when you're new to Go, so let's unpack it.

### Understanding slice expressions

Suppose our slice contains:

```text
Index:    0       1       2       3
Movie:    A       B       C       D
```

We want to delete movie C, which is at index `2`.

The expression:

```go
movies[:index]
```

means all elements before index `2`.

Result:

```text
[A, B]
```

The expression:

```go
movies[index+1:]
```

means all elements starting at index `3` and continuing to the end.

Result:

```text
[D]
```

Now we want to combine these two slices.

```go
append(movies[:index], movies[index+1:]...)
```

The `append()` function adds the elements of the second slice to the first.

But why is there `...`?

In Go, `append()` accepts individual elements after the first slice. The `...` expands a slice into its individual elements.

Conceptually, it becomes:

```go
append(movies[:index], movies[index+1], movies[index+2], ...)
```

So the final result is:

```text
[A, B, D]
```

We assign the resulting slice back to `movies`.

**Important:** this is a common Go slice-deletion idiom. It changes the slice's contents and length; it does not operate on a database or permanently delete data from disk.

After deleting a movie, our handler returns the remaining slice as JSON.

If the requested ID doesn't exist, it returns `404 Not Found`.

---

## 12. Routing: how requests reach the correct function

Defining handler functions is not enough. We also need to tell the server which function should handle each URL and HTTP method.

This is the job of the router.

Inside `main()`, we create one:

```go
r := mux.NewRouter()
```

Here, `:=` declares the variable `r` and initializes it with a new Gorilla Mux router.

Then we register our endpoints:

```go
r.HandleFunc("/movies", getMovies).Methods("GET")
r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
r.HandleFunc("/movies", createMovie).Methods("POST")
r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
```

Let's understand this syntax:

```go
r.HandleFunc("/movies", getMovies).Methods("GET")
```

- `r` is our router.
- `HandleFunc()` registers a URL pattern and its handler function.
- `"/movies"` is the URL path.
- `getMovies` is the function that should run when the route matches.
- `.Methods("GET")` restricts the route to GET requests.

The function name is passed without parentheses because we are registering the function, not calling it immediately.

### Why do we use `{id}`?

```go
"/movies/{id}"
```

The `{id}` part is a path variable.

It allows the same route pattern to match URLs such as:

```text
/movies/1
/movies/2
/movies/3
```

Mux extracts the actual value from the URL, making it available through `mux.Vars(r)`.

This means we don't need to create a separate route for every movie.

### Why can multiple routes use the same URL?

Notice that GET, PUT, and DELETE all use:

```text
/movies/{id}
```

This is intentional. The router considers both the path and the HTTP method when matching a request.

For example:

| Request | Handler |
|---|---|
| `GET /movies/3` | `getMovie` |
| `PUT /movies/3` | `updateMovie` |
| `DELETE /movies/3` | `deleteMovie` |

The URL identifies the movie, while the HTTP method determines the operation.

---

## 13. The `main()` function: starting everything

Our application starts in `main()`.

The overall structure looks like this:

```go
func main() {
    r := mux.NewRouter()

    // Initialize movies here.

    // Register routes here.

    fmt.Println("Starting server at port 8080")

    log.Fatal(http.ListenAndServe(":8080", r))
}
```

The actual implementation also initializes the six sample movies before registering the routes.

### Initializing the movies

We use `append()` to add our sample movie structs to the global `movies` slice.

This means the API has some data to work with immediately after starting.

### Starting the server

```go
log.Fatal(http.ListenAndServe(":8080", r))
```

This line starts the HTTP server.

Let's break it down:

**`http.ListenAndServe()`**

This function listens for incoming HTTP connections and serves requests using the supplied handler.

Its first argument is:

```go
":8080"
```

The colon followed by `8080` specifies the port on which the server listens. With no explicit host specified, the server listens on the available local interfaces.

Its second argument is:

```go
r
```

This is our Gorilla Mux router. It receives incoming requests and directs them to the matching handler.

**`log.Fatal()`**

If `ListenAndServe()` returns an error, `log.Fatal()` logs it and terminates the program.

A successfully running server normally keeps running inside `ListenAndServe()` until it encounters an error or is stopped.

That's why the terminal appears occupied when the application is running. This is expected.

---

## 14. Running the project locally

Make sure Go is installed and available in your terminal.

Open the VS Code terminal in the project folder and run:

```powershell
go mod init movie-api
go get github.com/gorilla/mux
go run .
```

If your project already has a `go.mod` file, don't run `go mod init` again. Just install the dependency if needed and run the project.

The first command initializes the Go module. The second downloads and registers Gorilla Mux as a dependency. The third builds and runs the application.

You should see:

```text
Starting server at port 8080
```

The server is now available at:

```text
http://localhost:8080
```

You can test the endpoints using Postman, another API client, or suitable command-line tools.

### Testing the endpoints

| Method | URL | What it does |
|---|---|---|
| GET | `http://localhost:8080/movies` | Retrieves every movie |
| GET | `http://localhost:8080/movies/1` | Retrieves movie 1 |
| POST | `http://localhost:8080/movies` | Creates a movie |
| PUT | `http://localhost:8080/movies/1` | Replaces movie 1 |
| DELETE | `http://localhost:8080/movies/1` | Deletes movie 1 |

For POST and PUT, select a raw JSON request body in Postman.

Example body:

```json
{
    "isbn": "100007",
    "title": "Interstellar",
    "director": {
        "firstname": "Christopher",
        "lastname": "Nolan"
    }
}
```

For POST, the server assigns the ID. For PUT, the URL determines which movie is updated, and the original ID is preserved.

---

## 15. A few Go concepts worth remembering

### `:=` versus `=`

These two operators look similar, but they do different things.

```go
id := "3"
```

`:=` declares a new local variable and assigns a value to it. Go infers the type.

```go
id = "4"
```

`=` assigns a new value to an already declared variable.

You generally cannot use `:=` when you only want to reassign an existing variable and declare no new variables.

### `var`

```go
var movie Movie
```

This explicitly declares a variable. If you don't provide an initial value, Go initializes it with the type's zero value.

### `range`

```go
for index, item := range movies {
    // Use index and item
}
```

`range` lets you iterate over collections such as slices and maps.

For slices, the first value is the index and the second value is the element.

If you don't need one of the values, use `_` to ignore it.

### Pointers and `&`

```go
json.NewDecoder(r.Body).Decode(&movie)
```

The `&` operator obtains the address of a variable.

The decoder needs to modify `movie`, so we give it a pointer to the variable rather than a copy of the variable's value.

The request parameter `r *http.Request` also uses a pointer type.

You don't need to master every detail of pointers to build this API, but you should understand that `&movie` lets the decoder write into the original variable.

### Multiple return values

Go functions can return more than one value.

For example:

```go
id, err := strconv.Atoi("123")
```

The first value is the converted integer, and the second is an error value.

This is why Go code frequently checks `err` after operations that may fail.

### `if`, `for`, and `return`

These are the control-flow tools used throughout the project.

- `if` executes a block only when its condition is true.
- `for` repeats a block or iterates over a collection.
- `return` exits the current function.

Together, they let us search through movies, detect invalid input, and stop execution when the desired result has been found.

---

## 16. What this project teaches about backend engineering

Even though the application is small, it introduces several important concepts.

**HTTP fundamentals:** You have a server that listens on a port, accepts requests, and sends responses.

**Routing:** Gorilla Mux maps URL patterns and HTTP methods to specific functions.

**Serialization:** JSON encoding converts Go structs and slices into data that clients can understand. JSON decoding converts incoming data into Go values.

**Data modeling:** Structs let you represent related data cleanly.

**CRUD operations:** You have implemented the basic operations used by many web services.

**Error handling:** Your handlers can return meaningful HTTP errors for malformed JSON and nonexistent movie IDs.

**Dependency management:** The project uses Go modules to manage an external library.

**State management:** You have seen how an in-memory collection can hold data while a program runs, and why persistent storage is needed for a real application.

These are fundamental building blocks. The same concepts appear in much larger applications, even when those applications use databases, authentication, middleware, and more complicated business logic.

---

## 17. Current limitations and what could come next

This is a learning project, not a production-ready service. A few limitations are worth keeping in mind.

1. **In-memory storage:** All changes disappear when the server restarts. A database such as PostgreSQL or SQLite would solve this.
2. **ID generation:** The current numeric ID strategy is simple, but not safe against simultaneous requests or every possible ID format.
3. **Validation:** The server checks whether JSON can be decoded, but it doesn't yet verify that required fields such as `title` or `isbn` are nonempty.
4. **Full replacement:** The PUT handler expects a complete movie representation. Sending only a title can erase the other fields.
5. **Concurrent access:** Multiple requests can access the global slice concurrently. A production implementation should protect shared mutable state with appropriate synchronization or use a storage layer that handles concurrency.
6. **HTTP status codes:** The creation handler returns `201 Created`, and missing movies return `404 Not Found`. The other handlers could be improved with more consistent status-code and error handling.
7. **Security and deployment:** The API has no authentication, rate limiting, or production deployment configuration. Those are separate concerns to learn when moving beyond a local exercise.

A sensible next step is to add input validation and tests. After that, you could move the storage layer to a database and learn how a more realistic Go backend is structured.

---

## Final takeaway

The biggest thing to remember is that this project is not just five functions that happen to work.

It's a complete, small request-response cycle:

1. The client sends an HTTP request.
2. Gorilla Mux matches the request to a route.
3. The appropriate handler receives the request and extracts any relevant data.
4. The handler reads or modifies the movie slice.
5. Go's JSON package serializes the result.
6. The HTTP response is sent back to the client.

Once that flow makes sense, backend development becomes much less mysterious. You can start with a simple in-memory slice, then gradually introduce databases, validation, middleware, authentication, testing, and deployment.

For now, understanding **why each line exists** is far more valuable than memorizing the code. That's what will let you build your own APIs instead of just following tutorials.
