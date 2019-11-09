package main

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/machinebox/graphql"
)

type rates struct {
	RateUSD float64
	RateEUR float64
}

type response struct {
	Items struct {
		Title string
	}
}

const (
	CurrencyUsd string = "USD"
	CurrencyEur string = "EUR"
)

var client *graphql.Client

var addr = "zenrus.ru:8888"

func init() {
	client = graphql.NewClient("http://localhost:8080/query")
}

func parse(data []byte) (rates, error) {
	str := string(data)
	RateUSD, err := strconv.ParseFloat(strings.Split(str, ";")[0], 64)
	RateEUR, err := strconv.ParseFloat(strings.Split(str, ";")[1], 64)
	return rates{RateEUR, RateUSD}, err
}

func startwebsocket(url string) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func gqlSend(currency string, rate float64) {
	req := graphql.NewRequest(`
		mutation($currency:Currency! $rate:Float!){
			updateRate(input: {
				currency: $currency
				exchangeRate: $rate
			})
		}
	`)
	req.Var("currency", currency)
	req.Var("rate", rate)

	ctx := context.Background()
	var res response
	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}
}

func main() {

	currentRate := rates{0, 0}

	ch := make(chan rates, 1)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c := startwebsocket(u.String())
	defer c.Close()

	go func(out chan<- rates) {
		defer close(out)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error trying to read message:", err)
				c = startwebsocket(u.String())
			}
			values, err := parse(message)
			if err != nil {
				log.Println("Error trying to parse message:", err)
				c = startwebsocket(u.String())
			}
			out <- values
		}
		return
	}(ch)

	for rate := range ch {
		if rate.RateEUR != currentRate.RateEUR {
			gqlSend(CurrencyEur, rate.RateEUR)
			currentRate.RateEUR = rate.RateEUR
		}
		if rate.RateUSD != currentRate.RateUSD {
			gqlSend(CurrencyUsd, rate.RateUSD)
			currentRate.RateEUR = rate.RateUSD
		}
	}
}
