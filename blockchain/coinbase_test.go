package blockchain

import (
	"fmt"
	"testing"
)

func TestCalcCoinbaseAmt(t *testing.T) {
	fmt.Println(CalcCoinbaseAmt(1))
	fmt.Println(CalcCoinbaseAmt(2016))
	fmt.Println(CalcCoinbaseAmt(2017))
	fmt.Println(CalcCoinbaseAmt(2016 * 2))
	fmt.Println(CalcCoinbaseAmt(2016 * 2+1))
	fmt.Println(CalcCoinbaseAmt(2016 * 3))
	fmt.Println(CalcCoinbaseAmt(2016 * 4))
	fmt.Println(CalcCoinbaseAmt(2016 * 5))
}
