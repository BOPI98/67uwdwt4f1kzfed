# LIBRARY_MANAGEMENT_SYSTEM 67uwdwt4f1kzfed

## Improvements

- Restructured the http server, the serveMux routing, the handlers into a separate package to use it in the integration test

- Added an integration test

- Added unit tests to the test package
  
- Fixed the unknown book/borrower error cases in the borrowers.GetBorrower and books.BorrowBook functions

- Added opentelemetry logger

- Added memory heap gauge and a request duration histogram metric

- Instrumented the http handler with otelhttp middleware

- I used the otelsql package to instrument the mysql connection

### Database setup
For the project I used 10.6.3-MariaDB database.

The Database.sql will create the database and the tables.
```
mysql -u root -p < Database.sql
```

### Golang:
The server was written with golang 1.24.4 

ini.v1 package was used for reading the config file

the mysql driver package is used for the database

opentelemetry packages was used for the metrics and traces

the go-sqlmock package was used for unit test database mocking 
```
go get -u gopkg.in/ini.v1
go get -u github.com/go-sql-driver/mysql

go get go.opentelemetry.io/otel/sdk/metric
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace
go get go.opentelemetry.io/otel/exporters/stdout/stdoutmetric

go get github.com/DATA-DOG/go-sqlmock
```

### Building
For building the server use the build.bat file that builds the executable to the build directory, then starts it.

### Config
The config file will store the database information, and the server port

### Tests
For running the test use the command below
```
go test .\test -v
```

### Tasks
Failed to write integration tests.

Incomplete unit test.

### REST API Endpoints

#### /books GET, POST

**GET**

request: -

response:

```
list:
    bookId      int
    title       string
    author      string
    pageCount   int
    borrowedBy  int | null
    createdAt   string (format: RFC3339)
```


**Status 200 OK**

Content type: application/json

```json
[
    {
        "bookId": 1,
        "title": "Book Title",
        "author": "Test Name",
        "pageCount": 200,
        "borrowedBy": null,
        "createdAt": "2024-01-20T10:00:00Z"
    }
]
```

**Status 400 Bad request**
Content type: text/plain
```
Error message
```

**POST**

request:

```
title       string
author      string
pageCount   int
borrowedBy  int | null
```

Content type: application/json

```json
{
    "title": "Book Title",
    "author": "Test Name",
    "pageCount": 200,
    "borrowedBy": null
}
```

response:

```
bookId      int (required)
```

**Status 200 OK**

Content type: application/json

```json
{
    "bookId": 12
}
```

**Status 400 Bad request**

Content type: text/plain

```
Error message
```

***

#### /books/{id}/borrow PATCH

**PATCH**

Path parameter

```
id: int
```

request:

```
borrowerId  int (required)
```

Content type: application/json

```json
{
    "borrowerId": 2
}
```

response:

**Status 200 OK**

*empty object*

Content type: application/json

```json
{}
```

**Status 400 Bad request**

Content type: text/plain

```
Error message
```

***

#### /borrowers POST

**POST**

request:


```
    fullname        string  (required)
    email           string  (required)
    phoneNumber     string  (required)
    age             int     (required)
```

Content type: application/json

```json
{
    "fullname": "Borrower Name",
    "email": "borrower@example.com",
    "phoneNumber": "0491570159",
    "age": 30
}
```

response:


**Status 200 OK**

*empty object*

Content type: application/json

```json
{}
```

**Status 400 Bad request**

Content type: text/plain

```
Error message
```

***

#### /borrowers/{id} GET

request: -

Path parameter

```
id: int
```

response:

```
fullname        string
email           string
phoneNumber     string
age             int
```

**Status 200 OK**

Content type: application/json

```json
{
    "fullname":     "Borrower Name",
    "email":        "borrower@example.com",
    "phoneNumber":  "0491570159",
    "age":          30
}
```

**Status 400 Bad request**

Content type: text/plain

```
Error message
```

***

#### /borrowers/{id}/borrowedbooks GET

Path parameter

```
id: int
```

response:

```
list:
    bookId      int
    title       string
    author      string
    pageCount   int
    borrowedBy  int | null
    createdAt   string (format: RFC3339)
```

**Status 200 OK**

Content type: application/json

```json
[
    {
        "bookId": 1,
        "title": "Book Title",
        "author": "Test Name",
        "pageCount": 200,
        "borrowedBy": null,
        "createdAt": "2024-01-20T10:00:00Z"
    }
]
```

**Status 400 Bad request**

Content type: text/plain

```
Error message
```
