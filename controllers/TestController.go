package controllers

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

var jsonList []map[string]string

func Test(c *gin.Context) {
	serialPath := "./conf/serial.txt"
	serialByte, _ := ioutil.ReadFile(serialPath)
	serial := string(serialByte)

	var html string
	jsonList = make([]map[string]string, 0)
	Bid := c.DefaultQuery("bid", serial)
	limit := c.DefaultQuery("limit", "10")
	intLimit, _ := strconv.Atoi(limit)

	for i := 0; i < intLimit; i++ {
		intBid, _ := strconv.Atoi(Bid)
		GetToken(intBid + i)
	}

	for _, v := range jsonList {
		title := gjson.Get(v["data"], "data.myDataSet.table.title")
		startDate := gjson.Get(v["data"], "data.myDataSet.table.startDate")
		content := gjson.Get(v["data"], "data.myDataSet.table.content")

		if title.Str != "" {
			intBid, _ := strconv.Atoi(v["Bid"])
			intSerial, _ := strconv.Atoi(serial)
			if intBid > intSerial {
				ioutil.WriteFile(serialPath, []byte(v["Bid"]), 0666)
			}
		}

		fmt.Printf("content is empty: %v\n", (content.Str == ""))

		html += "<h3>" + v["Bid"] + "<h3>"
		html += "<h5>startDate</h5>" + startDate.Str + "<br>"
		html += "<h5>title</h5>" + title.Str + "<br>"
		html += "<h5>data</h5>" + content.Str + "<br>"
		html += "<hr>"
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}

func GetToken(Bid int) {

	url := fmt.Sprintf("https://maplestory.beanfun.com/bulletin?bid=%v", Bid)

	crawler := colly.NewCollector() // 在colly中使用 Collector 這類物件 來做事情

	// 抓類別Class 名稱
	crawler.OnHTML("input[name='__RequestVerificationToken']", func(e *colly.HTMLElement) {
		token := e.Attr("value")
		cookie := e.Response.Headers.Get("Set-Cookie")
		GetData(strconv.Itoa(Bid), token, cookie)
	})

	crawler.OnRequest(func(r *colly.Request) { // iT邦幫忙需要寫這一段 User-Agent才給爬
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	})

	crawler.Visit(url)
}

func GetData(Bid string, token string, cookies string) {

	url := "https://maplestory.beanfun.com/bulletin?handler=BulletinDetail"
	requestData := map[string]string{
		"Bid": Bid,
	}

	crawler := colly.NewCollector() // 在colly中使用 Collector 這類物件 來做事情

	crawler.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", cookies)
		r.Headers.Set("X-CSRF-TOKEN", token)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	})

	crawler.OnResponse(func(r *colly.Response) {
		result := map[string]string{
			"Bid":  Bid,
			"data": string(r.Body),
		}
		jsonList = append(jsonList, result)
	})

	crawler.Post(url, requestData)
}
