FROM golang:1.24.6

# Current working directory
WORKDIR /go-to-s3

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy all the source code
COPY *.go ./

# Compile the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go-s3-client

EXPOSE 8081

# Tells docker what to run
CMD [ "/docker-go-s3-client"]