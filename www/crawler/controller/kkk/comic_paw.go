package kkk

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
	"time"
)

func cateList() (tags, regions, states map[string]int) {
	tags = map[string]int{
		"热血": 31,
		"恋爱": 26,
		"校园": 1,
		"冒险": 2,
		"后宫": 8,
		"科幻": 25,
		"战争": 12,
		"悬疑": 17,
		"推理": 33,
		"搞笑": 37,
		"奇幻": 14,
		"恐怖": 29,
		"神鬼": 20,
		"历史": 4,
		"同人": 30,
		"运动": 34,
		"绅士": 36,
		"机甲": 40,
	}
	regions = map[string]int{
		"港台": 35,
		"日韩": 36,
		"大陆": 37,
		"欧美": 52,
	}
	states = map[string]int{
		"连载中": 1,
		"已完结": 2,
	}

	return tags, regions, states
}

func ComicPaw() {
	tags, regions, states := cateList()
	for tag, tagId := range tags {
		for region, regionId := range regions {
			for state, stateId := range states {
				fmt.Printf("%s %s %s \n", tag, region, state)
				kk := common.Kind{
					Tag:    common.Kv{Name: tag, Val: tagId},
					Region: common.Kv{Name: region, Val: regionId},
					State:  common.Kv{Name: state, Val: stateId},
				}
				category(kk, "", 0)
			}
		}

	}
}

func ComicUpdate() {
	tags, regions, states := cateList()
	for tag, tagId := range tags {
		for region, regionId := range regions {
			for state, stateId := range states {
				kk := common.Kind{
					Tag:    common.Kv{Name: tag, Val: tagId},
					Region: common.Kv{Name: region, Val: regionId},
					State:  common.Kv{Name: state, Val: stateId},
				}
				category(kk, "-s18", 9)
			}
		}
	}
}

func category(kk common.Kind, sort string, limitPage int) {
	page := 1
	for {
		if limitPage != 0 && page > limitPage {
			break
		}
		url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/manhua-list-area%d-tag%d-st%d%s-p%d/",
			kk.Region.Val, kk.Tag.Val, kk.State.Val, sort, page)
		bot := robot.GetColly()

		existNode := false
		for i := 0; i <= 30; i++ {
			bot.OnHTML("div.mh-item", func(e *colly.HTMLElement) {
				existNode = true
				insertComic(e, kk)
			})

			err := bot.Visit(url)
			if err != nil {
				bot = robot.GetColly()
				continue
			} else {
				break
			}
		}

		if existNode == false {
			break
		}
		page++
	}
}

func insertComic(e *colly.HTMLElement, kk common.Kind) {
	title := e.DOM.Find(".title").Text()
	url, _ := e.DOM.Find(".title>a").Attr("href")
	id := tools.FindStringNumber(url)

	exists := new(model.SourceComic).Exists(id)
	if exists == true {
		return
	}
	sourceComic := new(model.SourceComic)
	sourceComic.Source = config.Spe.SourceId
	sourceComic.Title = title
	sourceComic.SourceId = id
	sourceComic.SourceUrl = "https://" + config.Spe.SourceUrl + url
	sourceComic.Label = model.Label{}
	sourceComic.LastChapterUpdateAt = time.Now().AddDate(-1, 0, 0)
	sourceComic.Category = kk.Tag.Name
	sourceComic.Region = kk.Region.Name
	if kk.State.Val == 2 {
		sourceComic.IsFinish = 1
	}
	err := orm.Eloquent.Create(&sourceComic).Error
	if err != nil {
		msg := fmt.Sprintf("漫画入库失败 source = %d source_id = %d err = %s", config.Spe.SourceId, id, err.Error())
		model.RecordFail(url, msg, "漫画入库", 1)
	} else {
		rd.RPush(common.SourceComicTASK, sourceComic.Id)
	}
}
