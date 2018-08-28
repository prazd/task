FROM golang:latest
RUN mkdir /gobot
ADD . /gobot
WORKDIR /gobot

RUN go get gopkg.in/tucnak/telebot.v2
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/anvie/port-scanner
RUN go get golang.org/x/crypto/bcrypt
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson

RUN go build -o main main.go

RUN go build -o mbot ./bot/bot.go

EXPOSE 80

CMD ["/gobot/mbot"]