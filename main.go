package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/melehin/gosplash/render"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/render.:format", func(c *gin.Context) {
		URL, ok := c.GetQuery("url")
		if !ok {
			c.String(http.StatusBadRequest, "url is not defined")
			return
		}
		URL, _ = url.QueryUnescape(URL)

		proxy, _ := c.GetQuery("proxy")
		viewport, _ := c.GetQuery("viewport")
		wait, _ := c.GetQuery("wait")
		timeout, _ := c.GetQuery("timeout")
		headless, _ := c.GetQuery("headless")
		images, _ := c.GetQuery("images")
		if viewport == "" {
			viewport = "1024x768"
		}
		format := c.Param("format")

		renderer, ok := render.Renderers[format]
		if !ok {
			c.String(http.StatusBadRequest, "renderer for %s is not found", format)
			return
		}

		contentType, data, err := render.Render(URL, proxy, viewport, wait, timeout, headless != "false" && headless != "0", images != "false" && images != "0", renderer)
		// error handling
		if err != nil {
			c.String(http.StatusBadGateway, "chromedp.Run error %v", err)
			return
		}
		if strings.Contains(string(data), "Chromium Authors") {
			c.String(http.StatusBadGateway, "%v", "Chrome show error page")
			return
		}
		if len(data) == 39 {
			c.String(http.StatusBadGateway, "%v", "Empty body in response")
			return
		}
		c.Data(http.StatusOK, contentType, data)
	})
	router.Run(":8050")
}
