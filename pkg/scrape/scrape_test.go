package scrape

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScrape(t *testing.T) {
	// articles := []Article{}
	// getWeaveWorksData("go", &articles, nil, nil)
	// fmt.Printf("%+v\n", articles)
	// require.Equal(t, "y", "x")
	str := "gogo go golang gogetem"
	count := strings.Count(str, " go ")
	require.Equal(t, 1, count)
}

func TestLimitArticles(t *testing.T) {
	article1 := Article{
		Title:      "article1",
		Occurences: 1,
	}
	article2 := Article{
		Title:      "article2",
		Occurences: 2,
	}
	article3 := Article{
		Title:      "article3",
		Occurences: 3,
	}
	articles := []Article{article3, article1, article2}
	articles = limitArticles(articles)
	expected := []Article{article3, article2, article1}
	require.Equal(t, expected, articles)
}
