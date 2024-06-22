package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oliver-hohn/screen_scraping_in_go/scrapers"
	"github.com/oliver-hohn/screen_scraping_in_go/scrapers/bbcgoodfood"

	_ "github.com/joho/godotenv/autoload"
)

var mode = flag.String("mode", "headed", "mode to run the scraper in: headed, headless, and remote")

func main() {
	flag.Parse()
	ctx := context.Background()

	var scraper *scrapers.Scraper
	switch *mode {
	case "headed":
		fmt.Printf("Running in headed mode\n")
		scraper = scrapers.NewHeadedScraper(ctx)
	case "headless":
		fmt.Printf("Running in headless mode\n")
		scraper = scrapers.NewHeadlessScraper(ctx)
	case "remote":
		fmt.Printf("Running in remote mode\n")
		scraper = scrapers.NewRemoteScraper(ctx, "http://localhost:9222")
	default:
		log.Fatalf("invalid mode: %s", *mode)
	}

	// Scrape recipes example
	recipes, err := bbcgoodfood.ScrapeRecipes(scraper)
	if err != nil {
		log.Fatal(err)
	}
	for _, recipe := range recipes {
		fmt.Printf("Recipe: %s\n", recipe.Name)
		fmt.Printf("Ingredients:\n")
		for _, ingredient := range recipe.Ingredients {
			fmt.Printf("  - %s\n", ingredient)
		}
		fmt.Printf("Link: %s\n", recipe.Link.String())
		fmt.Println()
	}

	// Post comment example
	if err := bbcgoodfood.PostComment(
		scraper,
		&bbcgoodfood.Credentials{
			Username: os.Getenv("BBC_GOOD_FOOD_USERNAME"),
			Password: os.Getenv("BBC_GOOD_FOOD_PASSWORD"),
		},
		recipes[0].Link,
		"This recipe is great!",
	); err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Minute)
}
