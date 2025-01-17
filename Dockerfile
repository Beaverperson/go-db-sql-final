FROM golang:1.21

WORKDIR /bd_app

COPY go.mod go.sum tracker.db ./

RUN go mod download 

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_bd_app 
# run with the test to verify sql base reachability
CMD ["/my_bd_app", "go test"]