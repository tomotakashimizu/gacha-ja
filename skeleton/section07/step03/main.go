// STEP03: ガチャを行うWebアプリを作ろう

package main

import (
	"fmt"
	"gacha-ja/skeleton/section07/step03/gacha"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func renderTemplate(w http.ResponseWriter, tmpl string, results []*gacha.Card) error {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		return err
	}
	return t.Execute(w, results)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// なぜか呼ばれない
	fmt.Println("サーバーが起動しました")
}

func run() error {
	p := gacha.NewPlayer(10, 100)
	// ※本当はハンドラ間で競合が起きるのでNG
	play := gacha.NewPlay(p)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		/* テンプレートに結果の一覧を埋め込んでレスポンスにする */
		if err := renderTemplate(w, "skeleton/section07/step03/index", play.Results()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// fmt.Fprintln(w, err.Error())
		}
	})

	// "ガチャを引く"submitボタンが押されると呼ばれる
	http.HandleFunc("/draw", func(w http.ResponseWriter, r *http.Request) {
		// r.FormValueメソッドを使ってフォームで入力したガチャの回数を取得
		// ガチャを行う回数は"num"で取得できる
		num, err := strconv.Atoi(r.FormValue("num"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i := 0; i < num; i++ {
			if !play.Draw() {
				break
			}
		}

		if err := play.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// "/"（トップ）にhttp.StatusFoundのステータスでリダイレクトする
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// ここは呼ばれる
	fmt.Println("サーバーが起動しました")
	return http.ListenAndServe(":8080", nil)
}
