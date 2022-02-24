package indicator

import sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"

func SMA(candles []sdk.Candle, length int) float64 {
	var summa float64
	for i := 0; i < length; i++ {
		summa += candles[i].ClosePrice
	}
	return summa / float64(length)
}

func EMA(candles []sdk.Candle, length int) float64 {
	if length == 1 {
		return candles[0].ClosePrice
	} else {
		alpha := 2 / float64(length+1)
		return alpha*candles[0].ClosePrice + (1-alpha)*calculatingEMA(candles, length-2, alpha, 0)
	}
}

func calculatingEMA(candles []sdk.Candle, length int, alpha float64, index int) float64 {
	index += 1
	if length == 0 {
		return candles[index].ClosePrice
	} else {
		return alpha*candles[index].ClosePrice + (1-alpha)*calculatingEMA(candles, length-1, alpha, index)
	}
}
