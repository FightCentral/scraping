package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/katana/pkg/engine/standard"
	"github.com/projectdiscovery/katana/pkg/output"
	"github.com/projectdiscovery/katana/pkg/types"
)

func main() {
	options := &types.Options{
		MaxDepth:     2,
		FieldScope:   "rdn",
		BodyReadSize: math.MaxInt,
		Timeout:      10,
		Concurrency:  10,
		Parallelism:  10,
		Delay:        5,
		RateLimit:    10,
		Strategy:     "breadth-first",
	}

	fighterURLs := []string{}
	urlMutex := &sync.Mutex{}
	maxFighters := 3

	options.OnResult = func(result output.Result) {
		url := result.Request.URL
		if strings.Contains(url, "/fighter-details") {
			urlMutex.Lock()
			fighterURLs = append(fighterURLs, url)
			urlMutex.Unlock()
		}
	}

	crawlerOptions, err := types.NewCrawlerOptions(options)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	defer crawlerOptions.Close()

	crawler, err := standard.New(crawlerOptions)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	defer crawler.Close()

	input := ""
	err = crawler.Crawl(input)
	if err != nil {
		gologger.Warning().Msgf("Could not crawl %s: %s", input, err.Error())
	}

	fighterURLs = uniqueStrings(fighterURLs)

	if len(fighterURLs) < maxFighters {
		maxFighters = len(fighterURLs)
		gologger.Warning().Msgf("Only %d fighter URLs found. Proceeding with available fighters.", maxFighters)
	} else {
		gologger.Info().Msgf("Limiting the trial run to %d fighters.", maxFighters)
	}

	c := colly.NewCollector(
		colly.AllowedDomains("", ""),
	)

	c.OnHTML("div.b-content", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("span.b-content__title-highlight"))
		stats := e.ChildText("ul.b-list__box-list li")
		fmt.Printf("Name: %s\n", name)
		for _, stat := range stats {
			fmt.Println(stat)
		}
		fmt.Println("---------------------------")
	})

	for _, url := range fighterURLs {
		err := c.Visit(url)
		if err != nil {
			log.Printf("Error visiting %s: %v", url, err)
		}
	}

	c.Wait()
}

func uniqueStrings(input []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
