package Chromedp

type JobData struct {
	Data apiData `json:"data"`
}

type apiData struct {
	Condition condition `json:"condition"`
	Welfare   welfare   `json:"welfare"`
	JobDetail jobDetail `json:"jobDetail"`
	Employees string    `json:"employees"`
}

type welfare struct {
	Welfare string `json:"welfare"`
}

type condition struct {
	Edu     string `json:"edu"`
	Other   string `json:"other"`
	WorkExp string `json:"workExp"`
}

type jobDetail struct {
	Description string `json:"jobDescription"`
	SalaryType  int32  `json:"salaryType"`
	SalaryMin   int32  `json:"salaryMin"`
	SalaryMax   int32  `json:"salaryMax"`
	Latitude    string `json:"Latitude"`
	Longitude   string `json:"Longitude"`
}
