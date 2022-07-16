package gacha

import "fmt"

type Player struct {
	Tickets int // ガチャ券の枚数
	Coin    int // コイン
}

// TODO: 引数にガチャ券とコインの枚数をもらい、
// それぞれをフィールドに設定したPlayer型の値を生成し、
// そのポインタを返すNewPlayer関数を作る
// func (p *Player) NewPlayer(tickets int, coin int) Player {
// 	p.Tickets = tickets
// 	p.Coin = coin
// 	return *p
// }

// メソッドをエクスポートする
// プレイヤーが行えるガチャの回数
func (p *Player) DrawableNum() int {
	// ガチャ券は1枚で1回、コインは10枚で1回ガチャが行える
	return p.Tickets + p.Coin/10
}

func (p *Player) draw(n int) {

	if p.DrawableNum() < n {
		fmt.Println("ガチャ券またはコインが不足しています")
		return
	}

	// ガチャ券から優先的に使う
	if p.Tickets >= n {
		p.Tickets -= n
		return
	}

	n -= p.Tickets
	p.Tickets = 0
	p.Coin -= n * 10 // 1回あたり10枚消費する
}
