// STEP04: ガチャの結果をデータベースに保存しよう

package main

import (
	"database/sql"
	"fmt"
	"gacha-ja/skeleton/section07/step04/gacha"
	"gacha-ja/skeleton/section07/step04/sqlite"
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
}

func run() error {
	fmt.Println("start run")

	db, err := sql.Open(sqlite.DriverName, "results.db")
	if err != nil {
		return fmt.Errorf("データベースのOpen:%w", err)
	}

	if err := createTable(db); err != nil {
		return err
	}

	p := gacha.NewPlayer(10, 100)
	// ※本当はハンドラ間で競合が起きるのでNG
	play := gacha.NewPlay(p)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET /")
		// データベースから結果を最大100件取得し、変数resultsに代入
		results, err := getResults(db, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := renderTemplate(w, "skeleton/section07/step04/index", results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// "ガチャを引く"submitボタンが押されると呼ばれる
	http.HandleFunc("/draw", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST /draw")

		num, err := strconv.Atoi(r.FormValue("num"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for i := 0; i < num; i++ {
			if !play.Draw() {
				break
			}

			if err := saveResult(db, play.Result()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if err := play.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	fmt.Println("サーバーが起動しました")
	return http.ListenAndServe(":8080", nil)
}

func createTable(db *sql.DB) error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS results(
		id        INTEGER PRIMARY KEY,
		rarity	  TEXT NOT NULL,
		name      TEXT NOT NULL
	);`

	_, err := db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("テーブル作成:%w", err)
	}

	return nil
}

func saveResult(db *sql.DB, card *gacha.Card) error {
	const sqlStr = `INSERT INTO results(rarity, name) VALUES (?,?);`
	// Execメソッドを用いてINSERT文を実行する
	_, err := db.Exec(sqlStr, card.Rarity.String(), card.Name)
	if err != nil {
		return err
	}
	return nil
}

func getResults(db *sql.DB, limit int) ([]*gacha.Card, error) {
	const sqlStr = `SELECT rarity, name FROM results LIMIT ?`
	rows, err := db.Query(sqlStr, limit)
	if err != nil {
		return nil, fmt.Errorf("%qの実行:%w", sqlStr, err)
	}
	defer rows.Close()

	var results []*gacha.Card
	for rows.Next() {
		var card gacha.Card
		// rows.Scanメソッドを用いてレコードをcardのフィールドに読み込む
		// cardに値が代入される
		err := rows.Scan(&card.Rarity, &card.Name)
		if err != nil {
			return nil, fmt.Errorf("Scan:%w", err)
		}
		results = append(results, &card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("結果の取得:%w", err)
	}

	return results, nil
}
