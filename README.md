# LIBRARY_MANAGEMENT_SYSTEM 67uwdwt4f1kzfed

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