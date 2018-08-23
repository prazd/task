package main

import (
	"os"
	"strconv"
	"time"

	"./mongo"
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

	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{volQ, invQ, volStatus, invStatus},
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Привет!Я помогу в мониторинге", &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
		b.Handle(&volQ, func(c *tb.Callback) {
			resp := mongo.QV()
			b.Edit(c.Message, "ℹ Vquantity: "+strconv.Itoa(resp), &tb.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

		b.Handle(&invQ, func(c *tb.Callback) {
			resp := mongo.QI()
			b.Edit(c.Message, "ℹ Iquantity: "+strconv.Itoa(resp), &tb.ReplyMarkup{
				InlineKeyboard: inlineKeys,
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
				resp += strconv.Itoa(i+1) + "." + arr[i][0] + " : " + arr[i][1]+"; Reviews    "+ arr[i][2] +"        "+ arr[i][3] + "\n"
			}

			b.Edit(c.Message, resp, &tb.ReplyMarkup{
				InlineKeyboard: inlineKeys,
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
				InlineKeyboard: inlineKeys,
			})
			b.Respond(c, &tb.CallbackResponse{})
		})

	})

	b.Start()
}
