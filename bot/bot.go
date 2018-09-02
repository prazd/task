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

	// Back
	BackToMain := tb.InlineButton{
		Unique: "BM",
		Text:   "Main",
	}

	BackToServices := tb.InlineButton{
		Unique: "BS",
		Text:   "🔙",
	}

	BackToDocker := tb.InlineButton{
		Unique: "BD",
		Text:   "🔙",
	}

	// Mongo
	MongoServices := tb.InlineButton{
		Unique: "MongServ",
		Text:   "MongoDB 🍃",
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
		Text:   "Server 🌐",
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

	DockerAlert := tb.InlineButton{
		Unique: "DA",
		Text:   "OK",
	}

	DockerCompose := tb.InlineButton{
		Unique: "DC",
		Text:   "docker-compose",
	}

	ComposeBuildUp := tb.InlineButton{
		Unique: "DBU",
		Text:   "Build\n&&\nUp",
	}

	ComposeStop := tb.InlineButton{
		Unique: "CStop",
		Text:   "Stop",
	}

	// ALL
	StopAllServices := tb.InlineButton{
		Unique: "ASS",
		Text:   "Stop 🍃🌐",
	}

	StartAllServices := tb.InlineButton{
		Unique: "SSS",
		Text:   "Start 🍃🌐",
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

	dockerAlert := [][]tb.InlineButton{
		[]tb.InlineButton{DockerAlert, BackToServices},
		[]tb.InlineButton{BackToMain},
	}

	dockerInline := [][]tb.InlineButton{
		[]tb.InlineButton{DockerStatus},
		[]tb.InlineButton{DockerStart, DockerStop, DockerPs},
		[]tb.InlineButton{DockerCompose},
		[]tb.InlineButton{BackToServices},
		[]tb.InlineButton{BackToMain},
	}

	dockerCompose := [][]tb.InlineButton{
		[]tb.InlineButton{ComposeBuildUp},
		[]tb.InlineButton{ComposeStop},
		[]tb.InlineButton{BackToDocker},
		[]tb.InlineButton{BackToMain},
	}

	mainInline := [][]tb.InlineButton{
		[]tb.InlineButton{volQ, invQ, volStatus, invStatus},
		[]tb.InlineButton{Services},
	}

	servicesInline := [][]tb.InlineButton{
		[]tb.InlineButton{ServicesStatus},
		[]tb.InlineButton{MongoServices, GoServer, Docker},
		[]tb.InlineButton{StopAllServices, StartAllServices},
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
				active := make(chan string)
				go func() {
					dockeractive := exec.Command("systemctl", "is-active", "docker")
					var stdout bytes.Buffer
					dockeractive.Stdout = &stdout
					err := dockeractive.Run()
					if err != nil {
						log.Println(err)
					}
					active <- stdout.String()
				}()

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

				resp := "1.Mongo:" + mongoS + "\n" + "2.Server:" + serverS + "\n" + "3.Docker:" + <-active

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

			b.Handle(&BackToDocker, func(c *tb.Callback) {
				b.Edit(c.Message, "Docker", &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			// Services
			// Mongo

			b.Handle(&MongoServices, func(c *tb.Callback) {
				ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)
				mongoPort := ps.IsOpen(27017)
				var info string
				if mongoPort == false {
					info = "Mongo: ✖"
				} else {
					info = "Mongo: ✔"
				}
				b.Edit(c.Message, info, &tb.ReplyMarkup{
					InlineKeyboard: mongoInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&MongoStop, func(c *tb.Callback) {
				resp := Systemctl("stop", "mongodb")
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: mongoInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&MongoStart, func(c *tb.Callback) {
				resp := Systemctl("start", "mongodb")
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: mongoInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&MongoReboot, func(c *tb.Callback) {
				resp := Systemctl("restart", "mongodb")
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: mongoInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			//Server

			b.Handle(&GoServer, func(c *tb.Callback) {
				ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)
				serverPort := ps.IsOpen(3000)
				var info string
				if serverPort == false {
					info = "Server: ✖"
				} else {
					info = "Server: ✔"
				}

				b.Edit(c.Message, info, &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&ServerStop, func(c *tb.Callback) {

				b.Edit(c.Message, "find the process id...", &tb.ReplyMarkup{
					InlineKeyboard: serverInline})
				info := make(chan string)

				ID := ServerProcessID()

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
					} else {
						info <- "nice"
					}
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

				// // Check status of process
				ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)
				serverPort := ps.IsOpen(3000)

				var resp string
				if serverPort == false {
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
					} else {
						info <- "Starting..."
					}

				}()

				b.Edit(c.Message, <-info, &tb.ReplyMarkup{
					InlineKeyboard: serverInline})

				ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)
				serverPort := ps.IsOpen(3000)
				var resp string
				if serverPort == true {
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
				b.Edit(c.Message, "Перед запуском контейнеров, убедитесь, что сервисы выключены", &tb.ReplyMarkup{
					InlineKeyboard: dockerAlert})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&DockerAlert, func(c *tb.Callback) {
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
				}()
				resp := <-info
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&DockerStop, func(c *tb.Callback) {
				resp := Systemctl("stop", "docker")
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&DockerStart, func(c *tb.Callback) {
				resp := Systemctl("start", "docker")
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			// docker-compose
			b.Handle(&DockerCompose, func(c *tb.Callback) {
				b.Edit(c.Message, "docker-compose", &tb.ReplyMarkup{
					InlineKeyboard: dockerCompose})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&ComposeBuildUp, func(c *tb.Callback) {
				status := make(chan string)
				go func() { // docker-compose build
					build := exec.Command("docker-compose", "build")
					err := build.Run()
					if err != nil {
						log.Println(err)
						status <- "bad"
					} else {
						status <- "nice"
					}
				}()

				up := <-status
				var resp string

				if up == "bad" {
					resp = "Not build"
				} else {
					resp = "Build"
				}

				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})

				go func() { // docker-compose up
					up := exec.Command("docker-compose", "up")
					err := up.Run()
					if err != nil {
						log.Println(err)
						status <- "bad"
					} else {
						status <- "nice"
					}
				}()

				up = <-status
				if up == "bad" {
					resp = "Not Up"
				} else {
					resp = "Up"
				}

				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			b.Handle(&ComposeStop, func(c *tb.Callback) {
				status := make(chan string)
				go func() { // docker-compose up
					up := exec.Command("docker-compose", "stop")
					err := up.Run()
					if err != nil {
						log.Println(err)
						status <- "bad"
					} else {
						status <- "nice"
					}
				}()

				var resp string
				info := <-status
				if info == "bad" {
					resp = "Exception"
				} else {
					resp = "Stop!"
				}

				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: dockerInline})
				b.Respond(c, &tb.CallbackResponse{})
			})

			//AllServices

			b.Handle(&StopAllServices, func(c *tb.Callback) {
				mongoStop := Systemctl("stop", "mongodb")
				// Server Stop
				info := make(chan string)
				go func() {
					ID := ServerProcessID()
					serverkillcmd := exec.Command("kill", "-9", ID)
					var out bytes.Buffer
					var stderr bytes.Buffer
					serverkillcmd.Stdout = &out
					serverkillcmd.Stderr = &stderr
					err := serverkillcmd.Run()
					if err != nil {
						log.Println(err, stderr.String())
					}
					ps := portscanner.NewPortScanner("localhost", 2*time.Second, 5)
					serverPort := ps.IsOpen(3000)
					if serverPort == false {
						info <- "Server stopped"
					} else {
						info <- "Server not stopped"
					}

				}()
				// Check status <- CHECKTHIS THING
				serverStop := <-info
				resp := "1.🍃:" + mongoStop + "\n" + "2.🌐:" + serverStop
				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: servicesInline})
				b.Respond(c, &tb.CallbackResponse{})

			})

			b.Handle(&StartAllServices, func(c *tb.Callback) {

				ps := portscanner.NewPortScanner("localhost", 5*time.Second, 5)
				mongoStart := Systemctl("start", "mongodb")
				var serverStart string

				serverS := exec.Command("./StartServer.sh")
				err := serverS.Run()
				if err != nil {
					log.Println(err)
				}

				var serverPort bool
				var count int
				for {
					if count > 20 {
						serverStart = "Problems"
						break
					}
					serverPort = ps.IsOpen(3000)
					if serverPort == true {
						serverStart = "ServerStart"
						break
					}
					count++
				}
				resp := "1.🍃:" + mongoStart + "\n" + "2.🌐" + serverStart

				b.Edit(c.Message, resp, &tb.ReplyMarkup{
					InlineKeyboard: servicesInline})
				b.Respond(c, &tb.CallbackResponse{})

			})

		} else {
			b.Send(m.Sender, "К сожалению вы не можете писать этому боту")
		}
	})

	b.Start()
}

func Systemctl(thing, service string) string {
	mongocmd := exec.Command("systemctl", thing, service)
	err := mongocmd.Run()
	var resp string
	if err != nil {
		log.Println(err)
		if thing == "start" {
			resp = "Not start !"
		} else if thing == "stop" {
			resp = "Not stopped !"
		} else if thing == "restart" {
			resp = "Not rebooted !"
		}
	} else {
		if thing == "start" {
			resp = "Started:✔"
		} else if thing == "stop" {
			resp = "Stopped:✖"
		} else if thing == "restart" {
			resp = "Rebooted:✔"
		}
	}
	return resp
}

func ServerProcessID() string {
	serverID := exec.Command("lsof", "-ti:3000")

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
	return ID
}
