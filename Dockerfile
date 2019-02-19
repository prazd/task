FROM golang:latest
RUN mkdir /goapp
ADD . /goapp
ENV CONN="mongodb://mongo:27017"
WORKDIR /goapp

RUN go get golang.org/x/crypto/bcrypt
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/prazd/task/server/handlers

RUN go build -o main ./server
RUN echo "export Check=100023 >> .bashrc" 

EXPOSE 3000
CMD ['/goapp/main']


