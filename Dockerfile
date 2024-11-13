# syntax=docker/dockerfile:1


FROM golang:latest AS builder

 WORKDIR /

# COPY go.mod go.sum ./

# RUN go mod download

# COPY *.go ./

# RUN CGO_ENABLED=0 go build  ./cmd/main.go 


# EXPOSE 8080

CMD ["./main"]