package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"
	indicator "tinkoff/indicator"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var token = flag.String("token", "", "your token")

func main() {
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	flag.Parse()
	rest()
}

func rest() {
	client := sdk.NewRestClient(*token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stocks, err := client.Stocks(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	// ограничения на получение часовых свеч только за 7 дней
	for _, stock := range stocks {
		from := time.Now().AddDate(0, 0, -7)
		to := time.Now()
		candles := getCandles(client, stock, from, to)
		for len(candles) <= 200 {
			to = from
			from = from.AddDate(0, 0, -7)
			candles = append(getCandles(client, stock, from, to), candles...)
		}
		for i, j := 0, len(candles)-1; i < j; i, j = i+1, j-1 {
			candles[i], candles[j] = candles[j], candles[i]
		}
		log.Printf("SMA 200 %v %v\n", stock.Ticker, indicator.SMA(candles, 200))
		log.Printf("EMA 200 %v %v\n", stock.Ticker, indicator.EMA(candles, 200))
	}
}

func getCandles(client *sdk.RestClient, stock sdk.Instrument, from, to time.Time) []sdk.Candle {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	candles, err := client.Candles(ctx, from, to, sdk.CandleInterval1Hour, stock.FIGI)
	if err != nil {
		log.Fatalln(err)
	}
	return candles
}
