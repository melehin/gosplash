package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/melehin/gosplash/render"

	"github.com/gin-gonic/gin"
)

func renderHandler(c *gin.Context) {
	URL, ok := c.GetQuery("url")
	if !ok {
		c.String(http.StatusBadRequest, "url is not defined")
		return
	}
	URL, _ = url.QueryUnescape(URL)

	proxy, _ := c.GetQuery("proxy")
	viewport, _ := c.GetQuery("viewport")
	js, _ := c.GetQuery("js")

	var script string
	if c.Request.Method == "POST" && c.ContentType() == "application/javascript" && js != "false" && js != "0" {
		raw, _ := c.GetRawData()
		script = string(raw)
	}

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

	contentType, data, err := render.Render(URL, proxy, viewport, script, wait, timeout, headless != "false" && headless != "0", images != "false" && images != "0", renderer)
	// error handling
	if err != nil {
		statusCode := http.StatusBadGateway
		if strings.Contains(err.Error(), "context deadline exceeded") {
			statusCode = http.StatusGatewayTimeout
		}
		c.String(statusCode, "chromedp.Run error %v", err)
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
}

func main() {
	router := gin.Default()

	router.GET("/render.:format", renderHandler)
	router.POST("/render.:format", renderHandler)

	router.Run(":8050")
}
