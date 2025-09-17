FROM golang:1.24.6

# current working directory
WORKDIR /go-sqlite

# Install build dependencies for CGO + SQLite
RUN apt-get update && apt-get install -y gcc libsqlite3-dev

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy Source Code
COPY *.go ./

# Compile
RUN CGO_ENABLED=1 GOOS=linux go build -o /go-sqlite

EXPOSE 8080


CMD ["./go-sqlite"]
