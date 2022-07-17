// STEP03: ガチャを行うWebアプリを作ろう

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	// "strconv"
	"gacha-ja/skeleton/section07/step03/gacha"
)

var tmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<head><title>ガチャ</title></head>
	<body>
		<form action="/draw">
			<label for="num">枚数</input>
			<input type="number" name="num" min="1" value="1">
			<input type="submit" value="ガチャを引く">
		</form>
		<h1>結果一覧</h1>
		<ol>{{range .}}
		<li>{{.}}</li>
		{{end}}</ol>
	</body>
</html>`))

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	p := gacha.NewPlayer(10, 100)
	// ※本当はハンドラ間で競合が起きるのでNG
	play := gacha.NewPlay(p)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if /* TODO: テンプレートに結果の一覧を埋め込んでレスポンスにする */ err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/draw", func(w http.ResponseWriter, r *http.Request) {
		// TODO: r.FormValueメソッドを使ってフォームで入力したガチャの回数を取得
		// ガチャを行う回数は"num"で取得できる

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

		// TODO: "/"（トップ）にhttp.StatusFoundのステータスでリダイレクトする
	})

	return http.ListenAndServe(":8080", nil)
}
