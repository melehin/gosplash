package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/melehin/gosplash/render"

	"github.com/gin-gonic/gin"
)

const CtxMaxPort = "MaxPort"

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

	var port int32 = -1
	if headless == "false" || headless == "0" {
		if p, ok := c.Get(CtxMaxPort); ok {
			port, _ = p.(int32)
		}
		log.Printf("Port: %d", port)
	}

	contentType, data, err := render.Render(URL, c.Request.Header.Get("Referer"), proxy, viewport, c.Request.Cookies(), script, wait, timeout, int(port), images != "false" && images != "0", renderer)
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

type Server struct {
	port int32
}

func (s *Server) Handler(c *gin.Context) {
	c.Set(CtxMaxPort, atomic.AddInt32(&s.port, 1))
	defer atomic.AddInt32(&s.port, -1)
	renderHandler(c)
}

func main() {
	router := gin.Default()

	var s Server

	router.GET("/render.:format", s.Handler)
	router.POST("/render.:format", s.Handler)

	router.Run(":8050")
}
