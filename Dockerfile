FROM golang:latest as builder
ENV GOPATH=/go
ENV GO111MODULE=on
WORKDIR ${GOPATH}/src/github.com/Shikugawa/nature-chart-slack
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o ./dist/nature-chart-slack -i ./main.go

FROM alpine:latest
RUN apk add --update --no-cache ca-certificates tzdata && update-ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/Shikugawa/nature-chart-slack/dist .
RUN chmod +x ./nature-chart-slack
EXPOSE 3000
ENTRYPOINT [ "/app/nature-chart-slack" ]