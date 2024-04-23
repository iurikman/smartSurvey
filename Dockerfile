FROM golang:1.22 as builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . /usr/src/app
RUN CGO_ENABLED=0 GOOS=linux go build -o smartsurvey cmd/service/main.go

FROM debian:stable-slim
COPY --from=builder . ./bin/smartsurvey
ENTRYPOINT ["./bin/smartsurvey"]