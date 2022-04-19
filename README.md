# 網路爬蟲
 chromedp + goquery 爬取 104 公司職缺 api 資料 
 並存入 mongoDB 
 
 
## models

| chromedp | github.com/chromedp/chromedp  |
| -------- | -------- |
| goquery  | github.com/PuerkitoBio/goquery | 

## 啟動專案
1. 本地需先啟動 mongo 服務 , 以及安裝 google chrome
2. cd ./cmd && go run .

### ~~docker 啟動~~
~~由於 chromedp/headless-shell imgae 尚不支援 arm 架構
故未進行測試~~
