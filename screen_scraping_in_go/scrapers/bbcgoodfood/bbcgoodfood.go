package bbcgoodfood

import (
	"context"
	"fmt"
	"net/url"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/oliver-hohn/screen_scraping_in_go/scrapers"
)

const BASE_URL = "https://www.bbcgoodfood.com"

type Credentials struct {
	Username string
	Password string
}

type Recipe struct {
	Name        string
	Ingredients []string

	Link *url.URL
}

func ScrapeRecipes(s *scrapers.Scraper) ([]*Recipe, error) {
	steps := []chromedp.Action{}
	steps = append(steps, navigateToRecipes()...)
	steps = append(steps, acceptCookies()...)

	var recipes []*Recipe
	steps = append(steps, readRecipes(&recipes)...)

	if err := s.Run(steps...); err != nil {
		return nil, fmt.Errorf("unable to scrape recipes: %w", err)
	}

	return recipes, nil
}

func navigateToRecipes() []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(BASE_URL + "/recipes/collection/easy-dinner-recipes"),
	}
}

func acceptCookies() []chromedp.Action {
	return []chromedp.Action{
		chromedp.Click(`//button[contains(text(), 'Accept All')]`),
	}
}

func readRecipes(recipes *[]*Recipe) []chromedp.Action {
	return []chromedp.Action{chromedp.ActionFunc(func(ctx context.Context) error {
		// Find all recipe articles
		var recipeArticleNodes []*cdp.Node
		if err := chromedp.Nodes(`//article[@data-item-type='recipe']`, &recipeArticleNodes).Do(ctx); err != nil {
			return err
		}

		// For each, find the path to the recipe
		recipeURLs := make([]*url.URL, len(recipeArticleNodes))
		for i, recipeArticleNode := range recipeArticleNodes {
			var recipePath string
			if err := chromedp.AttributeValue(recipeArticleNode.FullXPath()+"//a", "href", &recipePath, nil).Do(ctx); err != nil {
				return err
			}

			u, err := url.Parse(BASE_URL + recipePath)
			if err != nil {
				return err
			}

			recipeURLs[i] = u
		}

		for _, u := range recipeURLs {
			if err := chromedp.Navigate(u.String()).Do(ctx); err != nil {
				return err
			}

			var recipeName string
			if err := chromedp.Text(`//h1[contains(@class, 'heading')]`, &recipeName).Do(ctx); err != nil {
				return err
			}

			var ingredientNodes []*cdp.Node
			if err := chromedp.Nodes(`//h2[contains(text(), 'Ingredients')]/following-sibling::section//li`, &ingredientNodes).Do(ctx); err != nil {
				return err
			}

			ingredients := make([]string, len(ingredientNodes))
			for i, ingredientNode := range ingredientNodes {
				if err := chromedp.Text(ingredientNode.FullXPath(), &ingredients[i]).Do(ctx); err != nil {
					return err
				}
			}

			*recipes = append(*recipes, &Recipe{
				Name:        recipeName,
				Ingredients: ingredients,
				Link:        u,
			})
		}

		return nil
	})}
}

func PostComment(s *scrapers.Scraper, c *Credentials, recipeURL *url.URL, comment string) error {
	steps := []chromedp.Action{}
	steps = append(steps, navigateToRecipe(recipeURL)...)
	steps = append(steps, acceptCookies()...)
	steps = append(steps, logIn(c)...)
	steps = append(steps, writeComment(comment)...)

	if err := s.Run(steps...); err != nil {
		return fmt.Errorf("unable to post comment: %w", err)
	}

	return nil
}

func navigateToRecipe(recipeURL *url.URL) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(recipeURL.String()),
	}
}

func logIn(credentials *Credentials) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Click(`(//a[contains(text(), 'Sign in')])[1]`),
		chromedp.Click(`//button[contains(descendant::text(), 'Log in with your email')]`),
		chromedp.SendKeys(`//input[@name='email']`, credentials.Username),
		chromedp.SendKeys(`//input[@name='password']`, credentials.Password),
		chromedp.Click(`//button[contains(descendant::text(), 'Log in')]`),
	}
}

func writeComment(comment string) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Click(`//span[contains(text(), 'Comment')]`),
		chromedp.SendKeys(`//textarea`, comment),
	}
}
