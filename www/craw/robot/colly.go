package robot

import (
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/tidwall/gjson"
	"math/rand"
	"net/http"
	"time"
)

func GetColly() *colly.Collector {
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)
	return bot
}

func GetProxy() string {
	api := config.Spe.ProxyUrl
	var saveData []string
	proxy := ""
	cache := "proxy:" + config.Spe.SourceUrl
	cacheProxy := rd.Get(cache)
	if cacheProxy != "" {
		err := json.Unmarshal([]byte(cacheProxy), &saveData)
		if err != nil {
			panic(err)
		}
		if len(saveData) > 0 {
			randomIndex := rand.Intn(len(saveData))
			return saveData[randomIndex]
		}
	}
	for try := 0; try <= 30; try++ {
		content, code, _ := tools.HttpRequest(api,
			"GET", "", map[string]string{}, []*http.Cookie{})
		if code == 200 {
			res := gjson.Parse(content)
			for _, d := range res.Get("data").Array() {
				saveData = append(saveData, "http://"+d.Get("ip").String()+":"+d.Get("port").String())
			}
			jsonData, err := json.Marshal(saveData)
			if err != nil {
				panic(err)
			}
			rd.Set(cache, string(jsonData), time.Minute*20)
			break
		}
		if code != 200 {
			t := time.NewTicker(time.Second * 30)
			<-t.C
		}
		t := time.NewTicker(time.Second * 1)
		<-t.C
	}
	if len(saveData) > 0 {
		randomIndex := rand.Intn(len(saveData))
		proxy = saveData[randomIndex]
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
