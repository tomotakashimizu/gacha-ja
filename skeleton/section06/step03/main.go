// STEP03: エラー処理をまとめよう

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/tomotakashimizu/gacha"
)

var (
	flagCoin int
)

func init() {
	flag.IntVar(&flagCoin, "coin", 0, "コインの初期枚数")
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	tickets, err := initialTickets()
	if err != nil {
		return err
	}
	p := gacha.NewPlayer(tickets, flagCoin)
	play := gacha.NewPlay(p)

	n := inputN(p)
	for play.Draw() {
		if n <= 0 {
			break
		}
		fmt.Println(play.Result())
		n--
	}

	if err := play.Err(); err != nil {
		return fmt.Errorf("ガチャを%d回引く:%w", n, err)
	}

	if err := saveResults(play.Results()); err != nil {
		return err
	}

	if err := saveSummary(play.Summary()); err != nil {
		return err
	}

	return nil
}

func initialTickets() (int, error) {
	if flag.NArg() == 0 {
		return 0, errors.New("ガチャチケットの枚数を入力してください")
	}

	num, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		return 0, fmt.Errorf("ガチャチケット数のパース(%q):%w", flag.Arg(0), err)
	}

	return num, nil
}

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

func saveResults(results []*gacha.Card) (rerr error) {
	f, err := os.Create("results.txt")
	if err != nil {
		return fmt.Errorf("result.txtの作成:%w", err)
	}

	defer func() {
		if err := f.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("result.txtのクローズ:%w", err)
		}
	}()

	for _, result := range results {
		fmt.Fprintln(f, result)
	}

	return nil
}

func saveSummary(summary map[gacha.Rarity]int) (rerr error) {
	f, err := os.Create("summary.txt")
	if err != nil {
		return fmt.Errorf("summary.txtの作成:%w", err)
	}

	defer func() {
		if err := f.Close(); err != nil && rerr == nil {
			rerr = fmt.Errorf("summary.txtのクローズ:%w", err)
		}
	}()

	for rarity, count := range summary {
		fmt.Fprintf(f, "%s %d\n", rarity, count)
	}

	return nil
}
