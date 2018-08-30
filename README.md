# Backend and monitoring bot for IOS application

IOS app - https://github.com/Korobych/SocialApp

![Golang](https://i.pinimg.com/474x/17/af/bb/17afbb0a62d3c22dc447dde3871411cc--android-star-wars.jpg)

## GObot + GOserver + mongodb

### Clone this rep and 
```
$ git clone https://github.com/prazd/task.git
$ cd task
$ go get ... - install dependences: 
(
    golang.org/x/crypto/bcrypt
    gopkg.in/mgo.v2
    gopkg.in/mgo.v2/bson
    github.com/anvie/port-scanner
    gopkg.in/tucnak/telebot.v2
)
$ go build -o bot bot/bot.go
$ go build -o server main.go
$ export DevOne="Telegram Id of first developer"
$ export DevTwo="Telegram Id of second developer"
$ export BotToken="Token of Telegram bot"
```


![Docker](https://www.fullstackpython.com/img/logos/docker-wide.png)

### Docker 
```
$ docker-compose build
$ docker-compose up
```
