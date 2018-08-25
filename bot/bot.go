package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"./mongo"
	portscanner "github.com/anvie/port-scanner"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("BotToken"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return
	}

	volQ := tb.InlineButton{
		Unique: "VQ",
		Text:   "🏃",
	}

	invQ := tb.InlineButton{
		Unique: "IQ",
		Text:   "♿",
	}

	volStatus := tb.InlineButton{
		Unique: "VS",
		Text:   "🆗",
	}

	invStatus := tb.InlineButton{
		Unique: "IS",
		Text:   "🆘",
	}

	ServicesStatus := tb.InlineButton{
		Unique: "SS",
		Text:   "Status",
	}

	BotLog := tb.InlineButton{
		Unique: "ML",
		Text:   "B📄",
	}

	ServerLog := tb.InlineButton{
		Unique: "SL",
		Text:   "S📄",
	}

	MongoReboot := tb.InlineButton{
		Unique: "MR",
		Text:   "Reboot",
	}

	MongoStop := tb.InlineButton{
		Unique: "MStop",
		Text:   "Stop",
	}

	MongoStart := tb.InlineButton{
		Unique: "MStop",
		Text:   "Start",
	}

	Services := tb.InlineButton{
		Unique: "S",
		Text:   "Services",
	}

	BackToMain := tb.InlineButton{
		Unique: "BM",
		Text:   "Main",
	}

	BackToServices := tb.InlineButton{
		Unique: "BS",
		Text:   "Back",
	}

	MongoServices := tb.InlineButton{
		Unique: "BS",
		Text:   "MongoDB",
	}

	mongoInline := [][]tb.InlineButton{
		[]tb.InlineButton{MongoStart, MongoStop, MongoReboot, BackToServices, BackToMain},
	}

	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{volQ, invQ, volStatus, invStatus, Services},
	}

	servicesInline := [][]tb.InlineButton{
		[]tb.InlineButton{MongoServices, ServicesStatus, BotLog, ServerLog, BackToMain},
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Привет!Я помогу в мониторинге", &tb.ReplyMarkup{
			InlineKeyboard: mainInline,
		})
		b.Handle(&volQ, func(c *tb.Callback) {
			resp := mongo.QV()
			b.Edit(c.Message, "ℹ Vquantity: "+strconv.Itoa(resp), &tb.ReplyMarkup{
				InlineKeyboard: mainInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&invQ, func(c *tb.Callback) {
			resp := mongo.QI()
			b.Edit(c.Message, "ℹ Iquantity: "+strconv.Itoa(resp), &tb.ReplyMarkup{
				InlineKeyboard: mainInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&volStatus, func(c *tb.Callback) {
			arr := mongo.SV()
			var resp string
			for i, _ := range arr {
				if i == 0 {
					resp = "Готовы помочь:\n    🚹                   📱               👍     👎\n"
				}
				resp += strconv.Itoa(i+1) + "." + arr[i][0] + " : " + arr[i][1] + "; Reviews    " + arr[i][2] + "        " + arr[i][3] + "\n"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mainInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&invStatus, func(c *tb.Callback) {
			arr := mongo.SI()
			var resp string
			for i, _ := range arr {
				if i == 0 {
					resp = "Нужна помощь:\n       🆔                  🚹\n"
				}
				resp += strconv.Itoa(i+1) + "." + arr[i][0] + " : " + arr[i][1] + "\n"
			}
			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mainInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&ServicesStatus, func(c *tb.Callback) {
			ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)
			mongoPort := ps.IsOpen(27017)
			serverPort := ps.IsOpen(3000)
			var mongoS, serverS string
			if mongoPort == true {
				mongoS = "✔"
			} else {
				mongoS = "✖"
			}

			if serverPort == true {
				serverS = "✔"
			} else {
				serverS = "✖"
			}
			resp := "1.Mongo:" + mongoS + "\n" + "2.Server" + serverS

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: servicesInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})
		b.Handle(&BotLog, func(c *tb.Callback) {
			logfile, err := ioutil.ReadFile("bot.log")
			if err != nil {
				log.Println(err)
			}
			resp := string(logfile)

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: servicesInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&ServerLog, func(c *tb.Callback) {
			logfile, err := ioutil.ReadFile("main.log")
			if err != nil {
				log.Println(err)
			}
			resp := string(logfile)

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: servicesInline,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&Services, func(c *tb.Callback) {
			b.Edit(c.Message, "Services monitoring", &tb.ReplyMarkup{
				InlineKeyboard: servicesInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		// Back Handle

		b.Handle(&BackToMain, func(c *tb.Callback) {
			b.Edit(c.Message, "Monitoring info in mongodb", &tb.ReplyMarkup{
				InlineKeyboard: mainInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&BackToServices, func(c *tb.Callback) {
			b.Edit(c.Message, "Services", &tb.ReplyMarkup{
				InlineKeyboard: servicesInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		// Services
		// Mongo

		b.Handle(&MongoServices, func(c *tb.Callback) {
			b.Edit(c.Message, "Mongo", &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&MongoStop, func(c *tb.Callback) {
			mongocmd := exec.Command("systemctl", "stop", "mongodb")
			err := mongocmd.Run()
			resp := "Mongo: ✔"
			if err != nil {
				log.Println(err)
			} else {
				resp = "Mongo: ✖"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&MongoStart, func(c *tb.Callback) {
			mongocmd := exec.Command("systemctl", "start", "mongodb")
			err := mongocmd.Run()
			resp := "Mongo: ✖"
			if err != nil {
				log.Println(err)
			} else {
				resp = "Mongo: ✔"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&MongoReboot, func(c *tb.Callback) {
			mongocmd := exec.Command("systemctl", "restart", "mongodb")
			err := mongocmd.Run()
			resp := "Mongo: ✖"
			if err != nil {
				log.Println(err)
			} else {
				resp = "Mongo: ✔"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		//Server ....

	})

	b.Start()
}
