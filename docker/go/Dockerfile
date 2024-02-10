FROM golang:1.17-alpine3.16

COPY . /app

WORKDIR /app

RUN apk update \
&& apk add --no-cache git \
&& go get github.com/gin-gonic/gin@v1.7.4 \
&& go get github.com/jinzhu/gorm \
&& go get github.com/go-sql-driver/mysql

EXPOSE 8080
