package kk

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
	"github.com/tebeka/selenium"
	"path/filepath"
	"strings"
	"time"
)

func cateList() (tags, regions, pays, states map[string]int) {
	tags = map[string]int{
		"恋爱":  20,
		"古风":  46,
		"穿越":  80,
		"大女主": 77,
		"青春":  47,
		"非人类": 92,
		"奇幻":  22,
		"都市":  48,
		"总裁":  52,
		"强剧情": 82,
		"玄幻":  63,
		"系统":  86,
		"悬疑":  65,
		"末世":  91,
		"热血":  67,
		"萌系":  62,
		"搞笑":  71,
		"重生":  89,
		"异能":  68,
		"冒险":  93,
		"武侠":  85,
		"竞技":  72,
		"正能量": 54,
	}
	regions = map[string]int{
		"国漫": 2,
		"韩漫": 3,
		"日漫": 4,
	}
	pays = map[string]int{
		"全部": 0,
	}
	states = map[string]int{
		"连载中": 1,
		"已完结": 2,
	}

	return tags, regions, pays, states
}

func ComicPaw() {
	tags, regions, pays, states := cateList()
	for tag, tagId := range tags {
		for region, regionId := range regions {
			for pay, payId := range pays {
				for state, stateId := range states {
					fmt.Printf("%s %s %s %s \n", tag, region, pay, state)
					kk := common.Kind{
						Tag:    common.Kv{Name: tag, Val: tagId},
						Region: common.Kv{Name: region, Val: regionId},
						Pay:    common.Kv{Name: pay, Val: payId},
						State:  common.Kv{Name: state, Val: stateId},
					}
					category(kk, 2, 0)
				}
			}
		}
	}
}

func ComicUpdate() {
	tags, regions, pays, states := cateList()
	for tag, tagId := range tags {
		for region, regionId := range regions {
			for pay, payId := range pays {
				for state, stateId := range states {
					kk := common.Kind{
						Tag:    common.Kv{Name: tag, Val: tagId},
						Region: common.Kv{Name: region, Val: regionId},
						Pay:    common.Kv{Name: pay, Val: payId},
						State:  common.Kv{Name: state, Val: stateId},
					}
					category(kk, 3, 9)
				}
			}
		}
	}
}

func category(kk common.Kind, sort, limitPage int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/tag/%d?region=%d&pays=%d&state=%d&sort=%d&page=%d",
		kk.Tag.Val, kk.Region.Val, kk.Pay.Val, kk.State.Val, sort, 1)

	bot := robot.GetColly()
	totalPage := 1
	for try := 0; try <= 5; try++ {
		bot.OnHTML(".navigation", func(e *colly.HTMLElement) {
			last := e.DOM.Find(".itemBten").Last()
			totalPage = tools.FindStringNumber(last.Text())
		})

		err := bot.Visit(url)
		if err != nil {
			bot = robot.GetColly()
			if try > 1 {
				bot.SetProxy(robot.GetProxy())
			}
			if try == 5 {
				model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
				return
			}
		} else {
			break
		}
	}

	for page := 1; page <= totalPage; page++ {
		if limitPage > 0 && page > limitPage {
			continue
		}
		paw(kk, sort, page)
	}
}

func paw(kk common.Kind, sort, page int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/tag/%d?region=%d&pays=%d&state=%d&sort=%d&page=%d",
		kk.Tag.Val, kk.Region.Val, kk.Pay.Val, kk.State.Val, sort, page)
	A := ant.Get([]int{1})
	if A == nil {
		return
	}
	defer A.Free()
	A.Restart("")
	for try := 0; try <= 5; try++ {
		if try > 2 {
			A.Restart(robot.GetProxy())
		}
		A.WebDriver.Get(url)
		comicList, err := A.WebDriver.FindElements(selenium.ByClassName, "ItemSpecial")
		if err == nil {
			if try == 5 {
				model.RecordFail(url, "无法抓取分类列表页信息 :"+url, "列表错误", 0)
				return
			}
		}
		for _, comicElement := range comicList {
			state, _ := kk.State.Val.(int)
			insertComic(comicElement, kk.Tag.Name, kk.Region.Name, state)
		}
		break
	}
}

func insertComic(e selenium.WebElement, category string, region string, final int) {
	sourceComic := new(model.SourceComic)
	itemLink, err := e.FindElement(selenium.ByClassName, "itemLink")
	if err != nil {
		return
	}
	url, _ := itemLink.GetAttribute("href")
	sourceComic.SourceId = tools.FindStringNumber(url)
	if sourceComic.Exists(sourceComic.SourceId) {
		return
	}
	sourceComic.Source = config.Spe.SourceId
	sourceComic.Category = category
	sourceComic.Region = region
	sourceComic.Label = model.Label{}
	sourceComic.SourceUrl = url
	sourceComic.IsFinish = 0
	if final == 2 {
		sourceComic.IsFinish = 1
	}
	itemTitle, err := e.FindElement(selenium.ByClassName, "itemTitle")
	if err != nil {
		return
	}
	sourceComic.Title, _ = itemTitle.Text()
	if sourceComic.Title == "" {
		return
	}
	itemAuthor, err := e.FindElement(selenium.ByClassName, "author")
	if err == nil {
		sourceComic.Author, _ = itemAuthor.Text()
	}
	itemImg, err := e.FindElement(selenium.ByClassName, "img")
	if err == nil {
		cover0, _ := itemImg.GetAttribute("data-src")
		cover := strings.TrimSuffix(cover0, "-t.w207.webp.h")
		dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/kk_v_cover/%d/%d", config.Spe.SourceId, sourceComic.SourceId%128)
		for tryLimit := 0; tryLimit <= 7; tryLimit++ {
			proxy := ""
			if tryLimit > 3 {
				cover = cover0
			}
			if tryLimit > 5 {
				proxy = robot.GetProxy()
			}
			var cookies map[string]string
			cover := common.DownFile(cover, dir, filepath.Base(cover), proxy, cookies)
			if cover != "" {
				sourceComic.Cover = cover
				break
			}
		}
	}
	sourceComic.LastChapterUpdateAt = time.Now().AddDate(-1, 0, 0)
	err = orm.Eloquent.Create(&sourceComic).Error
	if err != nil {
		msg := fmt.Sprintf("漫画入库失败 source = %d source_id = %d err = %s", config.Spe.SourceId, sourceComic.SourceId, err.Error())
		model.RecordFail(url, msg, "漫画入库", 1)
	} else {
		rd.RPush(common.SourceComicTASK, sourceComic.Id)
	}
}
