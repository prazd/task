FROM golang:latest
RUN mkdir /gobot
ADD . /gobot
WORKDIR /gobot

RUN go get gopkg.in/tucnak/telebot.v2
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson

RUN go build bot.go

EXPOSE 80

CMD ["/gobot/bot"]

