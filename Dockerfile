FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app

# RUN apk update \
#     && apk --no-cache --update add build-base git

COPY ./fanfiction-backend/go.mod ./fanfiction-backend/go.sum ./

RUN go mod download && go mod tidy

COPY ./fanfiction-backend ./

RUN go build -o main ./cmd/api

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
CMD [ "/app/main" ]