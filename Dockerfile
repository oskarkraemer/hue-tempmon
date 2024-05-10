FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY templates ./templates

RUN go mod tidy

RUN go build -o /go-docker-huetemp


FROM alpine:latest
COPY --from=builder /go-docker-huetemp /go-docker-huetemp
COPY --from=builder /app/templates /templates

EXPOSE 8080
CMD [ "/go-docker-huetemp" ]