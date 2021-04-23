package scrape

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var weaveworksUrl = "https://www.weave.works/blog/"

const weaveworksLayout = "2006-01-02T15:04:05.999999999Z07:00"

func (scraper Scraper) getWeaveWorksData(topic string) {
	defer scraper.wg.Done()

	res, err := http.Get(weaveworksUrl)
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
	isFeaturedCard := true

	doc.Find(".card__body").EachWithBreak(func(index int, item *goquery.Selection) bool {
		if isFeaturedCard {
			isFeaturedCard = false
			return true
		}
		dateTag := item.Find("time")
		dateString, _ := dateTag.Attr("datetime")

		date, err := time.Parse(weaveworksLayout, dateString)
		if err != nil {
			fmt.Println(err)
			return true
		}
		if date.Before(time.Now().AddDate(0, 0, -7)) {
			return false
		}

		titleTag := item.Find(".card__title")
		title := titleTag.Text()
		linkTag := titleTag.Find("a")
		link, _ := linkTag.Attr("href")

		fullLink := fmt.Sprintf("https://www.weave.works%s", link)
		scraper.wg.Add(1)
		go scraper.createArticle("WeaveWorks", fullLink, topic, title)
		return true
	})
}
