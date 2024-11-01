package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/katana/pkg/engine/standard"
	"github.com/projectdiscovery/katana/pkg/output"
	"github.com/projectdiscovery/katana/pkg/types"
)

func main() {
	godotenv.Load()
	testMode := true

	var fighterURLs []string

	if testMode {
		fighterURLs = []string{
			"http://www.ufcstats.com/fighter-details/54f64b5e283b0ce7",
		}
	} else {
		fighterURLs = crawlFighterURLs()
	}

	save_fighters_data(fighterURLs)
	// var wg sync.WaitGroup
	// for _, url := range fighterURLs {
	// 	wg.Add(1)
	// 	go func(url string) {
	// 		defer wg.Done()
			
	// 		parseFighterURL(url)
	// 	}(url)
	// }
	// wg.Wait()
}

func crawlFighterURLs() []string {
	options := &types.Options{
		MaxDepth:     2,
		FieldScope:   "rdn",
		BodyReadSize: math.MaxInt,
		Timeout:      5,
		Concurrency:  5,
		Parallelism:  5,
		Delay:        1,
		RateLimit:    10,
		Strategy:     "breadth-first",
	}

	fighterURLs := []string{}
	urlMutex := &sync.Mutex{}

	options.OnResult = func(result output.Result) {
		url := result.Request.URL
		if strings.Contains(url, "/fighter-details/") {
			gologger.Info().Msgf("Found URL: %s", url)
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

	input := os.Getenv("URL") + "/statistics/fighters"
	err = crawler.Crawl(input)
	if err != nil {
		gologger.Warning().Msgf("Could not crawl %s: %s", input, err.Error())
	}

	fighterURLs = uniqueStrings(fighterURLs)

	fmt.Println("Found", len(fighterURLs), "fighter URLs")
	return fighterURLs
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
