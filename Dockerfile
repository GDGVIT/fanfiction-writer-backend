FROM golang:1.21.3-alpine3.18

WORKDIR /usr/src/app

# RUN apk update \
#     && apk --no-cache --update add build-base git

COPY ./fanfiction-backend/go.mod ./fanfiction-backend/go.sum ./

RUN go mod download && go mod tidy

COPY . .
