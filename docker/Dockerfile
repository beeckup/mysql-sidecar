FROM beeckup/base-image
ENV MINIO_ENABLED false
WORKDIR /go/src/app
COPY . .
RUN mkdir tempdump
RUN go build -o main src/main.go src/common.go src/minio.go src/aws.go src/mysql.go src/zip.go
CMD ["./main"]