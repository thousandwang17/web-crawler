// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package Chromedp

import (
	"Web-Crawler/internal/init/Connect"
	"Web-Crawler/internal/init/osenv"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// 需要獲取的資料
type data struct {
	Day         string // 徵才日期
	Bus         string // 徵才公司
	JobName     string // 徵才職稱
	Area        string // 徵才地區
	Experience  string // 徵才年資
	Education   string // 徵才教育程度
	Employees   string // 公司人員
	SalaryType  int32  // 月薪 日新 年薪
	SalarMin    int32  // 徵才薪資 Min
	SalarMax    int32  // 徵才薪資 Max
	GPSX        string
	GPSY        string
	Description string   // 徵才簡介
	Welfare     string   // 福利
	Tag         []string // 標籤
}

const (
	Request       = `https://www.104.com.tw/`
	JobContentApi = `/job/ajax/content`
)

// chromedp 上下文
var oz4ctx context.Context

// func checkChromePort() bool {
// 	addr := net.JoinHostPort("", "9222")
// 	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
// 	if err != nil {
// 		return false
// 	}
// 	defer conn.Close()
// 	return true
// }

// 初始化
func Search() {

	var ctx context.Context
	var cancel context.CancelFunc

	if osenv.TEST_MODE {
		ctx, _ = chromedp.NewContext(
			context.Background(),
		)
		ctx, cancel = chromedp.NewContext(ctx)
	} else {
		url := fmt.Sprintf("ws://%v:%v", osenv.CHORME_HANDLESS_ADDRESS, osenv.CHORME_HANDLESS_PORT)
		ctx, _ = chromedp.NewRemoteAllocator(context.Background(), url)
		ctx, cancel = chromedp.NewContext(ctx)
	}

	defer cancel()
	// 設定同樣的上下文
	oz4ctx = ctx

	getList()
}

// 搜尋工作清單
func getList() {
	for page := 1; page <= 3; page++ {
		var listHtml string

		//獲取 104 ul html
		err := chromedp.Run(oz4ctx,
			chromedp.Navigate(fmt.Sprintf("%s%v%s",
				`https://www.104.com.tw/jobs/search/?ro=0&kwop=7&keyword=golang&order=12&asc=0&page=`,
				page,
				`&mode=s&jobsource=2018indexpoc&langFlag=0`)),
			chromedp.Sleep(1*time.Second),
			chromedp.OuterHTML(`div#js-job-content`, &listHtml),
		)

		if err != nil {
			log.Fatalf(" err : %v", err)
		}

		// goquery將獲取到的 html 轉成 dom
		dom, err := goquery.NewDocumentFromReader(strings.NewReader(listHtml))

		if err != nil {
			log.Fatal(err)
		}

		// dom 抓取每筆 li
		var jobs []string
		dom.Find(`article.job-list-item`).Each(func(i int, selection *goquery.Selection) {
			html, err := selection.Html()
			if err != nil {
				log.Fatal(err)
			}
			jobs = append(jobs, html)
		})

		getJobsData(jobs)
	}
}

// 解析每筆 li 內的 html
func getJobsData(jobs []string) {

	mongo := Connect.GetMongo()
	collection := mongo.Database("web-crawler").Collection("Codeing-Jobs")

	defer func() {
		// 可以取得 panic 的回傳值
		r := recover()
		if r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	for _, v := range jobs {
		article, err := goquery.NewDocumentFromReader(strings.NewReader(v))

		if err != nil {
			log.Fatal(err)
		}

		// 徵才日期
		bus_data := data{}
		article.Find(`span.b-tit__date`).Each(func(i int, selection *goquery.Selection) {
			bus_data.Day = selection.Text()
		})

		// 徵才職缺名
		article.Find(`a.js-job-link`).Each(func(i int, selection *goquery.Selection) {
			bus_data.JobName = selection.Text()
		})

		// 徵才公司
		article.Find(`ul.b-list-inline > li > a`).Each(func(i int, selection *goquery.Selection) {
			bus_data.Bus = selection.Text()
		})

		// 徵才 標籤
		var tag []string
		article.Find(`div.job-list-tag > span`).Each(func(i int, selection *goquery.Selection) {
			tag = append(tag, selection.Text())
		})
		bus_data.Tag = tag

		// 徵才 地區, 年資, 教育程度
		var info []string
		article.Find(`ul.job-list-intro > li`).Each(func(i int, selection *goquery.Selection) {
			info = append(info, selection.Text())
		})

		if len(info) >= 3 {
			bus_data.Area = info[0]
			bus_data.Experience = info[1]
			bus_data.Education = info[2]
		}

		// 徵才詳細連結
		var jobDegitalLink string
		article.Find(`a.js-job-link`).Each(func(i int, selection *goquery.Selection) {
			if val, exists := selection.Attr("href"); exists {
				jobDegitalLink = val
			}
		})

		// 動態頁面 開新分頁 抓取 ajax
		if jobDegitalLink != `` {
			bus_data.getJobDegital(jobDegitalLink)
		}

		_, MGerr := collection.InsertOne(oz4ctx, bus_data)

		// _, MGerr := collection.UpdateOne(context.TODO(), filter, update, opts)

		if MGerr != nil {
			fmt.Println(MGerr)
		}

		empJSON, err := json.MarshalIndent(bus_data, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf(" %s\n", string(empJSON))
	}
}

// 動態頁面 監聽 ajax
func (jobData *data) getJobDegital(jobDegitalLink string) {

	// 開新分頁
	cloneCtx, cloneCancel := chromedp.NewContext(oz4ctx)

	defer cloneCancel()

	// JobContentApi  RequestID
	var jobRequestId string

	// 監聽開新分頁  Request
	chromedp.ListenTarget(cloneCtx, func(ev interface{}) {
		switch ev := ev.(type) {
		// 紀錄 JobContentApi 的 RequestID
		case *network.EventRequestWillBeSent:
			if ev.Type != "XHR" {
				return
			}

			if !strings.Contains(ev.Request.URL, JobContentApi) {
				return
			}

			jobRequestId = ev.RequestID.String()

		// 監聽 JobContentApi 的 Request 回傳
		case *network.EventResponseReceived:
			// 如果是 JobContentApi
			if ev.RequestID.String() == jobRequestId {
				// 取api ResponseBody
				go func() {
					c := chromedp.FromContext(cloneCtx)
					rbp := network.GetResponseBody(ev.RequestID)
					body, err := rbp.Do(cdp.WithExecutor(cloneCtx, c.Target))
					if err != nil {
						fmt.Println(err)
					}

					var jobDegital JobData
					if err_j := json.Unmarshal(body, &jobDegital); err_j != nil {
						log.Fatal(err_j)
					}
					jobData.Employees = jobDegital.Data.Employees
					jobData.SalaryType = jobDegital.Data.JobDetail.SalaryType
					jobData.SalarMin = jobDegital.Data.JobDetail.SalaryMin
					jobData.SalarMax = jobDegital.Data.JobDetail.SalaryMax
					jobData.GPSX = jobDegital.Data.JobDetail.Latitude
					jobData.GPSY = jobDegital.Data.JobDetail.Longitude
					// jobData.Description = jobDegital.Data.JobDetail.Description
					// jobData.Welfare = jobDegital.Data.Welfare.Welfare
				}()
			}
		}
	})
	// 開新分頁 瀏覽 jobDegitalLink
	err := chromedp.Run(cloneCtx,
		chromedp.Navigate(`https:`+jobDegitalLink),
		chromedp.Sleep(3*time.Second),
	)

	if err != nil {
		log.Fatal(err)
	}
}
