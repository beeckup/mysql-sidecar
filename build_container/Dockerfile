FROM golang:alpine
WORKDIR /go/src/app
COPY . .
RUN mkdir dumpdb
RUN apk add --no-cache git \
    && go get github.com/robfig/cron \
    && apk del git
RUN go build backup.go
RUN ls -la
CMD ["./backup"]