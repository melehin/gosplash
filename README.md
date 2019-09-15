# gosplash
Simple javascript renderer based on Chromium and Go with an HTTP API.
As a minimal alternative instead of https://github.com/scrapinghub/splash. A bit slower but no memory leaks.

# Run in container
```
wget https://raw.githubusercontent.com/jfrazelle/dotfiles/master/etc/docker/seccomp/chrome.json
docker run --rm -d -p 8050:8050 --security-opt seccomp=$(pwd)/chrome.json fedormelexin/gosplash
```
# HTTP API
## render.html
Return the HTML of the javascript-rendered page.
## render.png
Return the PNG of the javascript-rendered page.

## Arguments for both html/png:
**url : string : required** \
The url to render (required)

**headless : string : optional** \
headless=false used to prevent chrome headless detection. Chromium instance running under xvfb virtual desktop. \
Default headless=**true**

**timeout : string : optional** \
Examples: **30s** - 30 secs; **1m** - 1 minute \
A timeout for the render (defaults to 60s).

**wait : string : optional** \
Examples: **30s** - 30 secs; **1m** - 1 minute \
Time to wait for updates after page is loaded (defaults to 0).
> Wait time must be less than timeout.

**proxy : string : optional** \
A proxy URL

**viewport : string : optional** \
View width and height (in pixels) of the browser viewport to render the web page. Format is “<width>x<height>”, e.g. 800x600. Default value is 1024x768. 
‘viewport’ parameter is more important for PNG rendering; it is supported for all rendering endpoints because javascript code execution can depend on viewport size.
  
**images : integer : optional** \
Whether to download images. Possible values are 1 (download images) and 0 (don’t download images). 
> Default is 0 for html and 1 for png.

# Examples
(replace **localhost** with your external ip if necessary)
## Screenshot
```
wget -O screenshot.png http://localhost:8050/render.png?url=https://instagram.com/instagram&wait=1s&viewport=800x1280
```
## HTML page dump
```
wget -O index.html http://localhost:8050/render.html?url=https://instagram.com/instagram&wait=1s
```
