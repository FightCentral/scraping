package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Fighter struct {
	name     string
	nickname string
	record   string
	height   string
	weight   string
	reach    string
	stance   string
	dob      string
	SLpM     string
	StrAcc   string
	SApM     string
	StrDef   string
	TDAvg    string
	TDAcc    string
	TDDef    string
	SubAvg   string
}

func save_fighters_data(urls []string) {
	file, err := os.OpenFile("fighters.csv", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Error creating file to store fighters: %v", err)
	}
	defer file.Close()

	for _, url := range urls {
		fighter, err := parseFighterURL(url)

		if err != nil {
			log.Printf("Error parsing fighter URL %s: %v", url, err)
			continue
		}

		fmt.Println(fighter)
		result := fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", fighter.name, fighter.nickname, fighter.record, fighter.height, fighter.weight, fighter.reach, fighter.stance, fighter.dob, fighter.SLpM, fighter.StrAcc, fighter.SApM, fighter.StrDef, fighter.TDAvg, fighter.TDAcc, fighter.SubAvg)
		file.WriteString(result)
	}
}

func parseFighterURL(url string) (Fighter, error) {
	fighter := Fighter{}
	fmt.Println("Parsing fighter URL:", url)
	website := os.Getenv("URL")

	c := colly.NewCollector(
		colly.AllowedDomains(website),
	)

	c.OnHTML("section.b-statistics__section_details", func(e *colly.HTMLElement) {
		fighter.name = strings.TrimSpace(e.ChildText("h2.b-content__title > span.b-content__title-highlight"))
		fighter.nickname = strings.TrimSpace(e.ChildText("p.b-content__Nickname"))
		fighter.record = strings.TrimSpace(e.ChildText("h2.b-content__title > span.b-content__title-record"))

		// first value : record, second value: height, third value: weight, fourth value: reach, fifth value: stance, sixth value: dob
		e.ForEach("div.b-list__info-box_style_small-width ul.b-list__box-list li", func(i int, el *colly.HTMLElement) {
			title := strings.TrimSpace(el.ChildText("i.b-list__box-item-title"))
			value := strings.TrimSpace(el.DOM.Clone().ChildrenFiltered("i.b-list__box-item-title").Remove().End().Text())

			switch i {
			case 0:
				fighter.record = value
			case 1:
				fighter.height = value
			case 2:
				fighter.weight = value
			case 3:
				fighter.reach = value
			case 4:
				fighter.stance = value
			case 5:
				fighter.dob = value
			default:
				fmt.Println("Unknown stat:", title, value)
			}
		})

		// first value : SLpM, second value: Str. Acc., third value: SApM, fourth value: Str. Def
		e.ForEach("div.b-list__info-box-left > div.b-list__info-box-left > ul.b-list__box-list li", func(i int, el *colly.HTMLElement) {
			title := strings.TrimSpace(el.ChildText("i.b-list__box-item-title"))
			value := strings.TrimSpace(el.DOM.Clone().ChildrenFiltered("i.b-list__box-item-title").Remove().End().Text())

			switch i {
			case 0:
				fighter.SLpM = value
			case 1:
				fighter.StrAcc = value
			case 2:
				fighter.SApM = value
			case 3:
				fighter.StrDef = value
			default:
				fmt.Println("Unknown stat left:", title, value)
			}
		})
		
		// first value: ignore, second value: TDAvg, third value: TD Acc., fourth value: TD Def., fifth value: Sub. Avg.
		e.ForEach("div.b-list__info-box-left > div.b-list__info-box-right > ul.b-list__box-list li", func(i int, el *colly.HTMLElement) {
			title := strings.TrimSpace(el.ChildText("i.b-list__box-item-title"))
			value := strings.TrimSpace(el.DOM.Clone().ChildrenFiltered("i.b-list__box-item-title").Remove().End().Text())

			switch i {
			case 0:
			case 1:
				fighter.TDAvg = value
			case 2:
				fighter.TDAcc = value
			case 3:
				fighter.TDDef = value
			case 4:
				fighter.SubAvg = value
			default:
				fmt.Println("Unknown stat right:", title, value)
			}
		})
	})

	err := c.Visit(url)
	if err != nil {
		log.Printf("Error visiting %s: %v", url, err)
		return Fighter{}, err
	}

	c.Wait()

	fmt.Println(fighter)
	return fighter, nil
}
