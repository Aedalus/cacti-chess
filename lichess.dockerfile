FROM golang:1.16.4-buster AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/cacti-chess-uci ./uci
RUN go build -o /bin/cacti-chess-lichess ./lichess-bot

