package scrape

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var mediumUrl = "https://medium.com/tag"

const mediumLayout = "2006-01-02T15:04:05.999999999Z07:00"

func (scraper Scraper) getMediumData(topic string) {
	defer scraper.wg.Done()

	url := fmt.Sprintf("%v/%v/latest", mediumUrl, topic)
	res, err := http.Get(url)
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

	doc.Find(".postArticle.postArticle--short.js-postArticle.js-trackPostPresentation.js-trackPostScrolls").EachWithBreak(func(index int, item *goquery.Selection) bool {
		dateTag := item.Find("time")
		dateString, _ := dateTag.Attr("datetime")
		date, err := time.Parse(mediumLayout, dateString)
		if err != nil {
			fmt.Println(err)
			return true
		}

		if date.Before(time.Now().AddDate(0, 0, -7)) {
			return false
		}

		divTag := item.Find(".postArticle-readMore")
		linkTag := divTag.Find("a")
		link, _ := linkTag.Attr("href")
		titleTag := item.Find("h3")
		title := titleTag.Text()

		scraper.wg.Add(1)
		go scraper.createArticle("Medium", link, topic, title)
		return true
	})
}
