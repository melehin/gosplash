package render

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// DefaultTimeout sets the default number of seconds to timeout
const DefaultTimeout = 60 * time.Second

// Renderer sets the function interface for the renderer
type Renderer func(ctx context.Context) (string, []byte, error)

// Renderers contains a list of available renderers
var Renderers = map[string]Renderer{
	// get html dump
	"html": func(ctx context.Context) (string, []byte, error) {
		var dump string
		if err := chromedp.Run(ctx, chromedp.OuterHTML("html", &dump)); err != nil {
			return "", nil, fmt.Errorf("Could not get OuterHTML: %v", err)
		}
		return "text/html", []byte(dump), nil
	},
	// get screenshot
	"png": func(ctx context.Context) (string, []byte, error) {
		var buf []byte
		if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
			return "", nil, fmt.Errorf("Could not get screenshot: %v", err)
		}
		return "image/png", buf, nil
	},
}

// Render renders web page over Chromium instance
func Render(url, proxy, viewport, wait, timeout string, headless, images bool, r Renderer) (string, []byte, error) {
	// prepare options
	opts := chromedp.DefaultExecAllocatorOptions[:]

	opts = append(opts[:], chromedp.Flag("headless", headless))

	if proxy != "" {
		opts = append(opts[:], chromedp.ProxyServer(proxy))
	}

	if !images {
		opts = append(opts[:], chromedp.Flag("blink-settings", "imagesEnabled=false"))
	}

	vp := strings.Split(viewport, "x")
	if len(vp) == 2 {
		width, wok := strconv.Atoi(vp[0])
		height, hok := strconv.Atoi(vp[1])
		if wok == nil && hok == nil {
			opts = append(opts[:], chromedp.WindowSize(width, height))
		}
	}

	t := DefaultTimeout
	if tD, err := time.ParseDuration(timeout); err == nil {
		t = tD
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return "", nil, fmt.Errorf("could not navigate to %s: %v", url, err)
	}

	if w, err := time.ParseDuration(wait); err == nil {
		chromedp.Run(ctx, chromedp.Sleep(w))
	}

	return r(ctx)
}