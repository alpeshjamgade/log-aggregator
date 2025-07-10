FROM golang:1.23-alpine AS builder

RUN mkdir /app
COPY . /app
RUN ls
WORKDIR /app
RUN CGO_ENABLED=0 go build -o log-aggregator ./cmd/main.go
RUN chmod +x log-aggregator

# Install migrate CLI
RUN go install -ldflags="-s -w" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:latest
RUN mkdir /app

ENV ELASTIC_APM_RECORDING=false
ENV ELASTIC_APM_ACTIVE=false
ENV ELASTIC_APM_ENVIRONMENT=test
ENV ELASTIC_APM_SERVICE_NAME=log-aggregator

COPY --from=builder /app/log-aggregator /app
COPY --from=builder /go/bin/migrate /usr/local/bin/
COPY private /app/private

CMD [ "./app/log-aggregator" ]