package main

import (
	"Web-Crawler/init/Connect"
	"Web-Crawler/internal/Chromedp"
)

// 引入套件

// 程式執行入口
func main() {

	Connect.Mongo()

	// JobCrawler.Crawler()
	Chromedp.Search()
}
