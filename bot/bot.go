package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

	// Mongo
	MongoServices := tb.InlineButton{
		Unique: "MongServ",
		Text:   "MongoDB",
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
		Unique: "MStart",
		Text:   "Start",
	}

	// Server
	GoServer := tb.InlineButton{
		Unique: "GS",
		Text:   "Server",
	}

	ServerStart := tb.InlineButton{
		Unique: "ServStart",
		Text:   "Start",
	}

	ServerStop := tb.InlineButton{
		Unique: "ServStop",
		Text:   "Stop",
	}

	serverInline := [][]tb.InlineButton{
		[]tb.InlineButton{ServerStart, ServerStop, ServerLog, BackToServices, BackToMain},
	}

	mongoInline := [][]tb.InlineButton{
		[]tb.InlineButton{MongoStart, MongoStop, MongoReboot, BackToServices, BackToMain},
	}

	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{volQ, invQ, volStatus, invStatus, Services},
	}

	servicesInline := [][]tb.InlineButton{
		[]tb.InlineButton{MongoServices, GoServer, ServicesStatus, BotLog, BackToMain},
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
			resp := "Mongo: Not stopped ✔"
			if err != nil {
				log.Println(err)
			} else {
				resp = "Mongo: Stopped ✖"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&MongoStart, func(c *tb.Callback) {
			mongocmd := exec.Command("systemctl", "start", "mongodb")
			err := mongocmd.Run()
			resp := "Mongo: Not Start ✖"
			if err != nil {
				log.Println(err)
			} else {
				resp = "Mongo: Start ✔"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&MongoReboot, func(c *tb.Callback) {
			mongocmd := exec.Command("systemctl", "restart", "mongodb")
			err := mongocmd.Run()
			resp := "Mongo: not reboot ✖"
			if err != nil {
				log.Println(err)
			} else {
				resp = "Mongo: reboot ✔"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: mongoInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		//Server

		b.Handle(&GoServer, func(c *tb.Callback) {
			b.Edit(c.Message, "Server", &tb.ReplyMarkup{
				InlineKeyboard: serverInline})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&ServerStop, func(c *tb.Callback) {
			// Find id of processs
			serverID := exec.Command("lsof", "-t", "-i:3000")
			sOut, sErr := serverID.CombinedOutput()
			if sErr != nil {
				log.Println(sErr)
			}
			err := serverID.Run()
			if err != nil {
				log.Println(err)
			} else {
				b.Edit(c.Message, "Take process ID...", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
			}

			ID := string(sOut)
			ID = strings.Replace(ID, "\n", "", -1)

			// Stop Server
			serverkillcmd := exec.Command("kill", "-9", ID)
			_, sErr = serverkillcmd.CombinedOutput()
			if sErr != nil {
				log.Println(err)
			} else {
				b.Edit(c.Message, "Kill Server process...", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
			}

			err = serverkillcmd.Run()
			if err != nil {
				log.Println(err)
			}

			// Check status of process
			checkStatus := exec.Command("lsof", "-t", "-i:3000")
			sOut, sErr = serverID.CombinedOutput()
			if sErr != nil {
				log.Println(sErr)
			}
			err = checkStatus.Run()
			if err != nil {
				log.Println(err)
			}

			ID = string(sOut)
			var resp string

			if len(ID) == 0 {
				resp = "Server has stopped"
			} else {
				resp = "Server hasn't stopped"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: serverInline})
			b.Respond(c, &tb.CallbackResponse{})

		})

		b.Handle(&ServerStart, func(c *tb.Callback) {
			serverStart := exec.Command("./main", "2>", "server.log")
			err := serverStart.Run()
			if err != nil {
				log.Println(err)
			} else {
				b.Edit(c.Message, "Starting server...", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
			}

			checkStatus := exec.Command("lsof", "-t", "-i:3000")
			sOut, sErr := checkStatus.CombinedOutput()
			if sErr != nil {
				log.Println(sErr)
			}
			err = checkStatus.Run()
			if err != nil {
				log.Println(err)
			}
			if len(sOut) == 0 {
				b.Edit(c.Message, "Server hasn't start", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
			} else {
				b.Edit(c.Message, "Server has start", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
			}

			b.Respond(c, &tb.CallbackResponse{})
		})

	})

	b.Start()
}
