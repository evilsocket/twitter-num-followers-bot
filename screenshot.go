package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	profileBaseURL = "https://mobile.twitter.com/%s"
	waitFor        = "//div[@data-testid='primaryColumn']"
)

var (
	chrome = (*chromedp.CDP)(nil)
)

func setupChrome(ctx context.Context) (err error) {
	if chrome == nil {
		if chrome, err = chromedp.New(ctx); err != nil {
			return
		}
	}

	return
}

func takeScreenshot(ctx context.Context, screenName string, fileName string) (err error) {
	log.Printf("taking screenshot for user %s to file %s ...", screenName, fileName)

	if err = setupChrome(ctx); err != nil {
		return
	}

	var buf []byte
	err = chrome.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf(profileBaseURL, screenName)),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(waitFor, chromedp.BySearch),
		//chromedp.ScrollIntoView(`.banner-section.third-section`, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf("  writing %d bytes to %s", len(buf), fileName)
			return ioutil.WriteFile(fileName, buf, 0644)
		}),
	})

	if err == nil {
		err = chrome.Wait()
	}
	return
}
