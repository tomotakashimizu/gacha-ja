// STEP01: ファイルを分けよう

package main

import (
	"fmt"
	"gacha-ja/skeleton/section04/step02/gacha"
)

func main() {
	// ガチャ券10枚、コイン100枚を持ったプレイヤーを作る
	p := gacha.Player{Tickets: 10, Coin: 100}

	n := inputN(&p)
	// gacha.DrawN関数を呼び、変数resultsとsummaryに結果を代入する
	results, summary := gacha.DrawN(&p, n)

	fmt.Println(results)
	fmt.Println(summary)

	fmt.Println(p)
}

// 引数の型をgacha.Playerのポインタにする
func inputN(p *gacha.Player) int {

	max := p.DrawableNum()
	fmt.Printf("ガチャを引く回数を入力してください（最大:%d回）\n", max)

	var n int
	for {
		fmt.Print("ガチャを引く回数>")
		fmt.Scanln(&n)
		if n > 0 && n <= max {
			break
		}
		fmt.Printf("1以上%d以下の数を入力してください\n", max)
	}

	return n
}
