package main

import (
	"Web-Crawler/internal/Chromedp"
	"Web-Crawler/internal/init/Connect"

	_ "Web-Crawler/internal/init/osenv"
)

// 程式執行入口
func main() {
	Connect.Mongo()
	Chromedp.Search()
}
