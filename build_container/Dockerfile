FROM golang:alpine
WORKDIR /go/src/app
COPY . .
RUN mkdir dumpdb
RUN apk add --no-cache git \
    && go get github.com/robfig/cron \
    && apk del git
RUN go build backup.go
# install mysqldump
RUN apk update \
    && apk add mysql-client \
    && apk add --no-cache openssl curl
RUN apk add python
RUN apk add py-pip
RUN pip install awscli
RUN apk add --update coreutils
RUN chmod +x cron_script.sh
RUN chmod +x cleaner.sh
CMD ["./backup"]