package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/gocolly/colly"
	"math"
	"regexp"
	"strconv"
	"time"
)

func cateList() (tags, pays, states map[string]int) {
	tags = map[string]int{
		"恋爱": 105,
		"玄幻": 101,
		"异能": 103,
		"恐怖": 110,
		"剧情": 106,
		"科幻": 108,
		"悬疑": 112,
		"奇幻": 102,
		"冒险": 104,
		"犯罪": 111,
		"动作": 109,
		"日常": 113,
		"竞技": 114,
		"武侠": 115,
		"历史": 116,
		"战争": 117,
	}
	pays = map[string]int{
		"免费": 1,
		//"付费": 2,
	}
	states = map[string]int{
		"连载中": 1,
		"已完结": 2,
	}

	return tags, pays, states
}
func ComicPaw() {
	tags, pays, states := cateList()
	for tag, tagId := range tags {
		for pay, payId := range pays {
			for state, stateId := range states {
				tx := common.Kind{
					Tag:   common.Kv{Name: tag, Val: tagId},
					Pay:   common.Kv{Name: pay, Val: payId},
					State: common.Kv{Name: state, Val: stateId},
				}
				fmt.Printf("%s %s %s \n", tag, pay, state)
				category(tx, 0)
			}
		}
	}
}

func ComicUpdate() {
	tags, pays, states := cateList()
	for tag, tagId := range tags {
		for pay, payId := range pays {
			for state, stateId := range states {
				tx := common.Kind{
					Tag:   common.Kv{Name: tag, Val: tagId},
					Pay:   common.Kv{Name: pay, Val: payId},
					State: common.Kv{Name: state, Val: stateId},
				}
				category(tx, 7)
			}
		}
	}
}

func category(tx common.Kind, limitPage int) {
	bot := robot.GetColly()

	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/theme/%d/finish/%d/search/time/vip/%d/page/1",
		tx.Tag.Val, tx.State.Val, tx.Pay.Val)

	page := 1
	for try := 0; try <= 5; try++ {
		bot.OnResponse(func(r *colly.Response) {
			regexp := regexp.MustCompile(`var totalNum = "(\d+)";`)
			params := regexp.FindStringSubmatch(string(r.Body))
			if len(params) >= 2 {
				total, _ := strconv.Atoi(params[1])
				page = int(math.Ceil(float64(total) / float64(12)))
				for {
					if page < 1 {
						break
					}
					if limitPage > 0 && page > limitPage {
						page--
						continue
					}
					paw(tx, page)
					t := time.NewTicker(time.Second * 5)
					<-t.C
					page--
				}
			}
		})
		err := bot.Visit(url)
		if err != nil {
			bot = robot.GetColly()
			if try > 1 {
				bot.SetProxy(robot.GetProxy())
			}
			if try == 5 {
				model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
			}
		} else {
			break
		}
		t := time.NewTicker(time.Second * 2)
		<-t.C
	}
}

func paw(tx common.Kind, page int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/Comic/all/theme/%d/finish/%d/search/time/vip/%d/page/%d",
		tx.Tag.Val, tx.State.Val, tx.Pay.Val, page)
	bot := robot.GetColly()
	for try := 0; try <= 5; try++ {
		bot.OnHTML("li.ret-search-item", func(e *colly.HTMLElement) {
			state, _ := tx.State.Val.(int)
			insertComic(e, tx.Tag.Name, state)
		})

		err := bot.Visit(url)
		if err != nil {
			bot = robot.GetColly()
			if try > 1 {
				bot.SetProxy(robot.GetProxy())
			}
			if try == 5 {
				model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
			}
		} else {
			break
		}
	}
}

func insertComic(e *colly.HTMLElement, category string, final int) {
	info := e.DOM.Find(".ret-works-info")
	url, _ := info.Find(".ret-works-title>a").Attr("href")
	id := tools.FindStringNumber(url)
	author := info.Find(".ret-works-author").Text()
	coverUrl, _ := e.DOM.Find(".ret-works-cover img.lazy").Attr("data-original")
	popularity := e.DOM.Find(".ret-works-tags span").Last().Find("em").Text()

	exists := new(model.SourceComic).Exists(id)
	if exists == true {
		return
	}
	sourceComic := new(model.SourceComic)
	sourceComic.Source = config.Spe.SourceId
	sourceComic.SourceId = id
	sourceComic.SourceUrl = "https://" + config.Spe.SourceUrl + url
	sourceComic.Cover = coverUrl
	sourceComic.Author = author
	sourceComic.Label = model.Label{}
	e.ForEach(".ret-works-tags span", func(_ int, e *colly.HTMLElement) {
		text := e.DOM.Text()
		match, _ := regexp.MatchString("人气", text)
		if match == false {
			sourceComic.Label = append(sourceComic.Label, text)
		}
	})
	sourceComic.LastChapterUpdateAt = time.Now().AddDate(-1, 0, 0)
	sourceComic.Category = category
	sourceComic.Popularity = popularity
	if final == 2 {
		sourceComic.IsFinish = 1
	}
	var cookies map[string]string
	dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/%d/%d", config.Spe.SourceId, id%128)
	for tryLimit := 0; tryLimit <= 7; tryLimit++ {
		proxy := ""
		if tryLimit > 5 {
			proxy = robot.GetProxy()
		}
		cover := common.DownFile(sourceComic.Cover, dir, tools.RandStr(9)+".jpg", proxy, cookies)
		if cover != "" {
			sourceComic.Cover = cover
			break
		}
	}
	err := orm.Eloquent.Create(&sourceComic).Error
	if err != nil {
		msg := fmt.Sprintf("漫画入库失败 source = %d source_id = %d err = %s", config.Spe.SourceId, id, err.Error())
		model.RecordFail(url, msg, "漫画入库", 1)
	} else {
		rd.RPush(common.SourceComicTASK, sourceComic.Id)
	}
}
