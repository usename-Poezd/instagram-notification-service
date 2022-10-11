FROM golang:1.18-alpine3.14 AS build
WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app ./cmd/app
RUN go build -o /cron_job ./cmd/cron_job

######## Start a new stage from scratch #######
FROM alpine
WORKDIR /src

COPY --from=build /app  .
COPY --from=build /cron_job  .
RUN mkdir logs

EXPOSE 8000

CMD ["./app"]