package scrape

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var goblogUrl = "https://blog.golang.org/index"

const goblogLayout = "2 January 2006"

func (scraper Scraper) getGoBlogData(topic string) {
	defer scraper.wg.Done()

	res, err := http.Get(goblogUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find(".blogtitle").EachWithBreak(func(index int, item *goquery.Selection) bool {
		dateTag := item.Find(".date")
		date, err := time.Parse(goblogLayout, dateTag.Text())
		if err != nil {
			fmt.Println(err)
			return true
		}
		if date.Before(time.Now().AddDate(0, 0, -7)) {
			return false
		}

		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")
		title := linkTag.Text()

		fullLink := fmt.Sprintf("https://blog.golang.org%s", link)
		scraper.wg.Add(1)
		go scraper.createArticle("GoBlog", fullLink, topic, title)
		return true
	})
}
