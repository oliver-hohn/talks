package scrapers

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func Log(msg string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		fmt.Println(msg)
		return nil
	}
}

var _ chromedp.Action = Log("")
