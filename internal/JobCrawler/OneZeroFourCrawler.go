package JobCrawler

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

type data struct {
	Day        string
	Bus        string
	JobName    string
	Area       string
	Experience string
	Education  string
	Tag        []string
}

// 程式執行入口
func Crawler() {

	c := colly.NewCollector(
		colly.MaxDepth(2),
	) // 在colly中使用 Collector 這類物件 來做事情
	cDegital := c.Clone()

	cDegital.OnHTML("body", func(e *colly.HTMLElement) {
		fmt.Printf(" %v 456 \n", 123)

		fmt.Println(e)
	})

	c.OnHTML("article.job-list-item", func(e *colly.HTMLElement) {
		createDate := e.ChildText("span[class='b-tit__date']")
		jobName := e.ChildText("a[class='js-job-link']")
		bus := e.ChildText("ul.b-list-inline > li > a")
		info := e.ChildTexts("ul.job-list-intro > li")
		jobDegitalLink := e.ChildAttr("a.js-job-link", "href")

		tag := e.ChildTexts("div.job-list-tag > span")
		area := ""
		experience := ""
		education := ""

		if len(info) >= 3 {
			area = info[0]
			experience = info[1]
			education = info[2]
		}

		if createDate != "" {
			catchDate := data{
				Day:        createDate,
				Bus:        jobName,
				JobName:    bus,
				Area:       area,
				Experience: experience,
				Education:  education,
				Tag:        tag,
			}
			// Print link

			fmt.Printf(" %v %v \n", catchDate, "https:"+jobDegitalLink)
			time.Sleep(3 * time.Second)
			errorMsg := cDegital.Visit("https:" + jobDegitalLink)
			fmt.Printf(" %v 123 \n", errorMsg)

		}

	})

	c.OnRequest(func(r *colly.Request) { // iT邦幫忙需要寫這一段 User-Agent才給爬
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	})

	for page := 1; page <= 1; page++ {
		time.Sleep(2 * time.Second)
		err := c.Visit(
			fmt.Sprintf("%s%v%s",
				"https://www.104.com.tw/jobs/search/?ro=0&kwop=7&keyword=golang&expansionType=area%2Cspec%2Ccom%2Cjob%2Cwf%2Cwktm&order=12&asc=0&page=",
				page,
				"&mode=s&jobsource=2018indexpoc&langFlag=0"))

		if err != nil {
			fmt.Printf("err %v", err)
		}

	}

}
