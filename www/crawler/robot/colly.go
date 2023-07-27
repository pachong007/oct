package robot

import (
	"comics/tools"
	"comics/tools/config"
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
	content, code, _ := tools.HttpRequest("https://dvapi.doveproxy.net/cmapi.php?rq=distribute&user=yipinbao6688&token=eUkxbHhCSFZFcit1TS9XRWdxVy9mUT09&auth=0&geo=PH&city=208622&agreement=1&timeout=35&num=1&rtype=0",
		"GET", "", map[string]string{}, []*http.Cookie{})
	proxy := ""
	if code == 200 {
		res := gjson.Parse(content)
		proxy = "http://" + res.Get("data").Get("ip").String() + ":" + res.Get("data").Get("port").String()
	}
	t := time.NewTicker(time.Second * 1)
	<-t.C
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
