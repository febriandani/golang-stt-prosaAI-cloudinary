FROM golang:alpine

RUN apk update && apk add --no-cache git

#add curl
RUN apk --no-cache add curl

# SET TZ
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

WORKDIR /app

COPY . .

RUN go mod tidy -go=1.16 && go mod tidy -go=1.17

RUN go mod vendor

RUN go build -o binary /app/cmd/api/main.go

ENTRYPOINT ["/app/binary"]

EXPOSE 8004