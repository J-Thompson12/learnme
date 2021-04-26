package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/justin/learnme/pkg/scrape"
	"github.com/nlopes/slack"
)

var site = "https://blog.golang.org/index"

func main() {
	err := godotenv.Load("environment.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.HandleFunc("/learn", slashCommandHandler).Methods("Post")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}

func slashCommandHandler(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/learnme":
		w.Write([]byte("processing...\n"))
		params := &slack.Msg{Text: s.Text}
		go handleArticlesResponse(s.ResponseURL, params.Text)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleArticlesResponse(responseURL, topic string) {
	articles := scrape.Scrape(site, topic)
	reader := buildResponse(topic, articles)
	responseBody := bytes.NewBuffer(reader)
	resp, err := http.Post(responseURL, "application/json", responseBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

}
func buildResponse(topic string, articles []scrape.Article) []byte {
	response := ""
	if len(articles) == 0 {
		response = response + "Unable to find any matching articles"
	}

	for i, article := range articles {
		response = response + fmt.Sprintf("%v: %v\n     %v\n     %v\n", i+1, article.Site, article.Title, article.Link)
	}
	response = fmt.Sprintf("{'text': '%v'}", response)
	return []byte(response)
}
