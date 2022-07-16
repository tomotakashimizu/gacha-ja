package main

// rarityやcardに関する処理をここに移す
// -- card.goに移す ここから --
type rarity string

const (
	rarityN  rarity = "N"
	rarityR  rarity = "R"
	raritySR rarity = "SR"
	rarityXR rarity = "XR"
)

func (r rarity) String() string {
	return string(r)
}

type card struct {
	rarity rarity // レア度
	name   string // 名前
}

func (c *card) String() string {
	// レア度:名前のように文字列を作る
	// 例："SR:ドラゴン"
	return c.rarity.String() + ":" + c.name
}

// -- card.goに移す ここまで --
