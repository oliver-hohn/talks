package scrapers

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

type Scraper struct {
	cancelFuns []context.CancelFunc
	chromeCtx  context.Context
}

func NewHeadlessScraper(ctx context.Context) *Scraper {
	cancelFuncs := []context.CancelFunc{}

	opts := chromedp.DefaultExecAllocatorOptions[:]

	allowCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	cancelFuncs = append(cancelFuncs, cancel)

	chromeCtx, cancel := chromedp.NewContext(allowCtx, chromedp.WithErrorf(log.Printf))
	cancelFuncs = append(cancelFuncs, cancel)

	return &Scraper{
		cancelFuns: cancelFuncs,
		chromeCtx:  chromeCtx,
	}
}

func NewHeadedScraper(ctx context.Context) *Scraper {
	cancelFuncs := []context.CancelFunc{}

	// Override default "headless" options to run Chrome in headed mode (with a UI)
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts,
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
	)
	allocatorCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	cancelFuncs = append(cancelFuncs, cancel)

	// Initialize a new Chrome instance with the overriden options
	chromeCtx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithErrorf(log.Printf))
	cancelFuncs = append(cancelFuncs, cancel)

	return &Scraper{
		cancelFuns: cancelFuncs,
		chromeCtx:  chromeCtx,
	}
}

func (s *Scraper) Cancel() {
	for _, cancel := range s.cancelFuns {
		cancel()
	}
}

func (s *Scraper) Run(actions ...chromedp.Action) error {
	return chromedp.Run(s.chromeCtx, actions...)
}
