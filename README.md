
# Library Management System - Part I


## Project Overview

The service manages operations on books, such as searching,
adding, buying, removing, borrowing and returning books,
as well as adding copies to existing books.

Each book is identified by its unique title and contains information about its author, publication year, total number of copies and number of available copies.

In Part I, clients communicate with the service through a REST API, which uses gRPC calls to a server that returns dummy responses.

## Prerequisites

Before building and running the project, make sure the following tools are installed:

- Go >= 1.26
- Protocol Buffers compiler (protoc)
- GNU Make
- curl
The project was developed and tested on Ubuntu using WSL2.

The required code generators are installed using the make setup command.

## Build

Run the following commands from the project root directory.

Install the required code generators (one-time setup):

```sh
make setup
```

Generate the REST and gRPC code and build the project:

```sh
make build
```

The build command generates the required code from the OpenAPI and
protobuf specifications, checks the Go code using go vet, and compiles
the server.

The executable is created at `bin/server`.

## Run the Servers

Start the first server in a terminal:

```sh
make run
```

Start the second server in another terminal:

```sh
make run2
```

| Server | HTTP port | gRPC port |
|--------|-----------|-----------|
| n1     | 8080      | 50051     |
| n2     | 8081      | 50052     |

The first server handles REST requests by issuing gRPC calls
to its own gRPC service.

The second server handles REST requests by issuing gRPC calls
to the gRPC service of n1.

Both servers must remain running during the demonstration.
To stop a server, press Ctrl+C in its terminal.


## Demonstration

The `etc/demo.sh` script demonstrates all eight REST operations,
including successful responses and error cases.

With both servers running, open a new terminal in the project
root directory.

To test the first server (n1), run:

```sh
make demo
```

To test the second server (n2), run:

```sh
./etc/demo.sh http://localhost:8081
```

The script displays the HTTP status codes and, for successful
operations with a response body, the returned data.

The responses from n2 contain the "from n1" marker in the returned
book data, demonstrating that n2 forwards RPCs to n1.

Since Part I uses dummy responses, write operations do not
permanently change the state of the library.


## REST API Overview

The Library Management System provides eight REST operations:

| Method | Path | Operation |
|---|---|---|
| GET | `/v1/AllBooks` | GetAllBooks|
| GET | `/v1/AllBooks/{Booktitle}` | GetBook |
| POST | `/v1/AllBooks` | AddBook |
| PATCH | `/v1/AllBooks/{Booktitle}` | AddCopies |
| DELETE | `/v1/AllBooks/{Booktitle}` | RemoveBook |
| POST | `/v1/AllBooks/{Booktitle}/borrow` | BorrowBook |
| POST | `/v1/AllBooks/{Booktitle}/return` | ReturnBook |
| POST | `/v1/AllBooks/{Booktitle}/buy` | BuyBook |

Read operations allow clients to choose between sequential
and linearizable consistency.

Write operations are defined as linearizable.

These consistency guarantees will be enforced in Part III.

The complete API specification is available in `api/openapi.yaml`.


## API Documentation

Once the server is running, the REST API can be viewed and
tested using the following tools:

- Swagger UI: http://localhost:8080/docs
- Swagger Editor: http://localhost:8080/editor
- OpenAPI specification: http://localhost:8080/openapi.yaml

Swagger UI allows users to explore the eight REST operations
and send HTTP requests using the "Try it out" feature.

Swagger Editor allows users to view, edit and validate the
OpenAPI specification.

The same tools are also available on the second server
through port 8081.

The REST API is defined in `api/openapi.yaml`, and the gRPC
service is defined in `proto/booklibrary/v1/booklibrary.proto`.


## Project Structure

| File or Directory | Description |
|---|---|
| `api/openapi.yaml` | REST API specification |
| `api/oapi-codegen.yaml` | REST code generator configuration |
| `api/spec.go` | Serves the OpenAPI specification |
| `api/docs.html` | Swagger UI |
| `api/editor.html` | Swagger Editor |
| `proto/booklibrary/v1/booklibrary.proto` | gRPC service and message definitions |
| `internal/grpcsvc/server.go` | gRPC service with dummy responses |
| `internal/httpapi/handlers.go` | REST handlers, RPC calls and error mapping |
| `cmd/server/main.go` | Initializes and connects the HTTP and gRPC servers |
| `etc/demo.sh` | Demonstration script |
| `Makefile` | Commands for generating, building and running the project |
| `go.mod` | Go module and dependencies |
| `gen/` | Automatically generated REST and gRPC code |

The generated code in `gen/` should not be edited manually.
It is regenerated using `make build`.

