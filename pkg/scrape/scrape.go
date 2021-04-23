package scrape

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Site       string
	Title      string
	Link       string
	Occurences int
}

type Scraper struct {
	wg       *sync.WaitGroup
	mux      *sync.Mutex
	articles *[]Article
}

func Scrape(website string, topic string) []Article {
	var wg sync.WaitGroup
	scraper := Scraper{
		wg:       &wg,
		mux:      &sync.Mutex{},
		articles: &[]Article{},
	}

	scraper.wg.Add(2)
	go scraper.getGoBlogData(topic)
	go scraper.getWeaveWorksData(topic)
	go scraper.getMediumData(topic)
	scraper.wg.Wait()

	return limitArticles(*scraper.articles)
}

func linkTopicOccurences(url, topic string) int {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Lower string and add spaces so count doesnt include the word as a substring of another word (go, golang)
	topic = strings.ToLower(topic)
	topic = fmt.Sprintf(" %v ", topic)

	return strings.Count(strings.ToLower(doc.Text()), strings.ToLower(topic))
}

func limitArticles(articles []Article) []Article {
	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].Occurences > articles[j].Occurences
	})
	if len(articles) > 10 {
		articles = articles[:10]
	}
	return articles
}

func (scraper Scraper) createArticle(site, link, topic, title string) {
	defer scraper.wg.Done()
	occurences := linkTopicOccurences(link, topic)
	if occurences == 0 {
		return
	}

	article := Article{
		Site:       site,
		Title:      title,
		Link:       link,
		Occurences: occurences,
	}
	scraper.mux.Lock()
	*scraper.articles = append(*scraper.articles, article)
	scraper.mux.Unlock()
}
