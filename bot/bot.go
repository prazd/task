package main

import (
	"bytes"
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

	// Docker
	Docker := tb.InlineButton{
		Unique: "Docker",
		Text:   "Docker 🐳",
	}

	DockerPs := tb.InlineButton{
		Unique: "DockerPS",
		Text:   "PS",
	}

	DockerStart := tb.InlineButton{
		Unique: "DockerStart",
		Text:   "Start",
	}

	DockerStop := tb.InlineButton{
		Unique: "DockerStop",
		Text:   "Stop",
	}

	DockerStatus := tb.InlineButton{
		Unique: "DockerStatus",
		Text:   "Status",
	}

	// Inline
	serverInline := [][]tb.InlineButton{
		[]tb.InlineButton{ServerStart, ServerStop},
		[]tb.InlineButton{ServerLog},
		[]tb.InlineButton{BackToServices},
		[]tb.InlineButton{BackToMain},
	}

	mongoInline := [][]tb.InlineButton{
		[]tb.InlineButton{MongoStart, MongoStop, MongoReboot},
		[]tb.InlineButton{BackToServices},
		[]tb.InlineButton{BackToMain},
	}

	dockerInline := [][]tb.InlineButton{
		[]tb.InlineButton{DockerStart, DockerStop},
		[]tb.InlineButton{DockerPs},
		[]tb.InlineButton{DockerStatus},
		[]tb.InlineButton{BackToServices},
		[]tb.InlineButton{BackToMain},
	}

	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{volQ, invQ, volStatus, invStatus},
		[]tb.InlineButton{Services},
	}

	servicesInline := [][]tb.InlineButton{
		[]tb.InlineButton{MongoServices, GoServer, ServicesStatus},
		[]tb.InlineButton{Docker},
		[]tb.InlineButton{BotLog},
		[]tb.InlineButton{BackToMain},
	}

	b.Handle("/start", func(m *tb.Message) {
		whitelist := [2]string{os.Getenv("DevOne"), os.Getenv("DevTwo")} // id's
		if strconv.Itoa(m.Sender.ID) == whitelist[0] || strconv.Itoa(m.Sender.ID) == whitelist[1] {
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

				b.Edit(c.Message, "find the process id...", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
				info := make(chan string)

				go func() {
					serverID := exec.Command("lsof", "-t", "-i:3000")
					sOut, sErr := serverID.CombinedOutput()
					if sErr != nil {
						log.Println(sErr)
					}
					ID := string(sOut)
					ID = strings.Replace(ID, "\n", "", -1)

					info <- ID

				}()

				ID := <-info
				// Stop Server

				b.Edit(c.Message, "kill the process...", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})

				go func() {
					serverkillcmd := exec.Command("kill", "-9", ID)
					var out bytes.Buffer
					var stderr bytes.Buffer
					serverkillcmd.Stdout = &out
					serverkillcmd.Stderr = &stderr
					err := serverkillcmd.Run()
					if err != nil {
						log.Println(err, stderr.String())
						info <- "bad"
					}
					info <- "nice"
				}()

				kill := <-info

				if kill == "bad" {
					b.Edit(c.Message, "Not killed", &tb.ReplyMarkup{
						InlineKeyboard: serverInline})
				} else {
					b.Edit(c.Message, "Killed", &tb.ReplyMarkup{
						InlineKeyboard: serverInline})
				}

				b.Edit(c.Message, "Check status of process....", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})

				// Check status of process
				go func() {
					serverID := exec.Command("lsof", "-t", "-i:3000")

					var stderr bytes.Buffer
					var stdout bytes.Buffer

					serverID.Stderr = &stderr
					serverID.Stdout = &stdout
					err := serverID.Run()
					if err != nil {
						log.Println(stderr.String(), err)
					}

					ID := stdout.String()
					ID = strings.Replace(ID, "\n", "", -1)
					info <- ID
				}()
				ID = <-info
				var resp string

				if len(ID) == 0 {
					resp = "Server stopped"
				} else {
					resp = "Server didn't stop"
				}

				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
				b.Respond(c, &tb.CallbackResponse{})

			})

			b.Handle(&ServerStart, func(c *tb.Callback) {
				info := make(chan string)
				go func() {
					serverStart := exec.Command("./StartServer.sh")
					err := serverStart.Run()
					if err != nil {
						log.Println(err)
						info <- "Fail"
					}
					info <- "Starting..."
				}()

				b.Edit(c.Message, <-info, &tb.ReplyMarkup{
					InlineKeyboard: serverInline})

				go func() {
					serverID := exec.Command("lsof", "-t", "-i:3000")

					var stderr bytes.Buffer
					var stdout bytes.Buffer

					serverID.Stderr = &stderr
					serverID.Stdout = &stdout
					err := serverID.Run()
					if err != nil {
						log.Println(stderr.String(), err)
					}

					ID := stdout.String()
					ID = strings.Replace(ID, "\n", "", -1)
					info <- ID
				}()

				var resp string
				if len(<-info) != 0 {
					resp = "Nice"
				} else {
					resp = "Fail"
				}

				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
				b.Respond(c, &tb.CallbackResponse{})

			})
			// Docker
			b.Handle(&Docker, func(c *tb.Callback) {
				b.Edit(c.Message, "Docker", &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&DockerPs, func(c *tb.Callback) {
				info := make(chan string)
				go func() {
					ps := exec.Command("docker", "ps")
					var stdout bytes.Buffer
					ps.Stdout = &stdout
					err := ps.Run()
					if err != nil {
						log.Println(err)
					}
					info <- stdout.String()
				}()
				b.Edit(c.Message, <-info, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&DockerStatus, func(c *tb.Callback) {
				info := make(chan string)
				go func() {
					ps := exec.Command("systemctl", "is-active", "docker")
					var stdout bytes.Buffer
					ps.Stdout = &stdout
					err := ps.Run()
					if err != nil {
						log.Println(err)
					}
					info <- stdout.String()
					b.Send(m.Sender, stdout.String())
				}()
				var resp string

				if <-info == "active" {
					resp = "Active"
				} else {
					resp = "Inactive"
				}
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&DockerStop, func(c *tb.Callback) {
				info := make(chan string)
				go func() {
					dockerstop := exec.Command("systemctl", "stop", "docker")
					err := dockerstop.Run()
					if err != nil {
						log.Println(err)
						info <- "bad"
					}
				}()
				if <-info == "bad" {
					b.Edit(c.Message, "Active", &tb.ReplyMarkup{
						InlineKeyboard: dockerInline})
					b.Respond(c, &tb.CallbackResponse{})
				} else {
					b.Edit(c.Message, "Inactive", &tb.ReplyMarkup{
						InlineKeyboard: dockerInline})
					b.Respond(c, &tb.CallbackResponse{})
				}
			})
			b.Handle(&DockerStart, func(c *tb.Callback) {
				info := make(chan string)
				go func() {
					dockerstop := exec.Command("systemctl", "start", "docker")
					err := dockerstop.Run()
					if err != nil {
						log.Println(err)
						info <- "bad"
					}
				}()
				var resp string
				if <-info != "bad" {
					resp = "Active"
				} else {
					resp = "Inactive"
				}
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})

			})

		} else {
			b.Send(m.Sender, "К сожалению вы не можете писать этому боту")
		}
	})

	b.Start()
}
