package scrapers

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func Log(msg string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		fmt.Println(msg)
		return nil
	}
}

var _ chromedp.Action = Log("")

func Measure(action chromedp.Action) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		start := time.Now()

		if err := action.Do(ctx); err != nil {
			return err
		}

		end := time.Now()
		fmt.Printf("Action took %s\n", end.Sub(start))

		return nil
	}
}

var _ chromedp.Action = Measure(nil)
