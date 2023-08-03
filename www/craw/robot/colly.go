package robot

import (
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/tidwall/gjson"
	"net/http"
	"time"
)

func GetColly() *colly.Collector {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)
	proxy := GetProxy()
	if proxy != "" && config.Spe.AppDebug == false {
		bot.SetProxy(proxy)
	}
	return bot
}

func GetProxy() string {
	proxy := ""
	cache := "proxy:" + config.Spe.SourceUrl
	cacheProxy := rd.Get(cache)
	if cacheProxy != "" {
		return cacheProxy
	}
	for {
		content, code, _ := tools.HttpRequest("https://dvapi.doveproxy.net/cmapi.php?rq=distribute&user=yipinbao6688&token=eUkxbHhCSFZFcit1TS9XRWdxVy9mUT09&auth=0&geo=NZ&city=361574&agreement=1&timeout=20&num=1&rtype=0",
			"GET", "", map[string]string{}, []*http.Cookie{})
		if code == 200 {
			res := gjson.Parse(content)
			proxy = "http://" + res.Get("data").Get("ip").String() + ":" + res.Get("data").Get("port").String()
			rd.Set(cache, proxy, time.Minute*15)
			break
		}
		if code == 409 {
			t := time.NewTicker(time.Second * 30)
			<-t.C
		}
		t := time.NewTicker(time.Second * 1)
		<-t.C
	}
	return proxy
}

func GetSeleniumArgs() map[string]string {
	if config.Spe.AppDebug {
		return map[string]string{
			"--user-agent": config.Spe.UserAgent,
			"--headless":   "-",
			"--no-sandbox": "-",
		}
	} else {
		return map[string]string{
			"--user-agent":            config.Spe.UserAgent,
			"--proxy-server":          GetProxy(),
			"--disable-dev-shm-usage": "",
		}
	}
}
