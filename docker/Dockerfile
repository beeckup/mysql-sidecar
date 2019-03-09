FROM golang:alpine
WORKDIR /go/src/app
COPY . .
RUN mkdir tempdump
RUN apk add --no-cache git mysql-client openssl curl python py-pip coreutils\
    && go get github.com/robfig/cron \
    && go get github.com/aws/aws-sdk-go \
    && go get github.com/go-sql-driver/mysql
RUN go build -o main main.go common.go minio.go aws.go mysql.go zip.go
CMD ["./main"]