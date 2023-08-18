package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ChapterPaw() {
	taskLimit := 20

	for limit := 0; limit < taskLimit; limit++ {
		common.StopSignal("章节任务挂起")
		id, err := rd.LPop(common.SourceComicTASK)
		if err != nil || id == "" {
			id, err = rd.LPop(common.SourceComicRenewTASK)
			if err != nil || id == "" {
				return
			}
		}
		for i := 0; i <= 10; i++ {
			bot := robot.GetColly()
			sourceComic := new(model.SourceComic)
			if orm.Eloquent.Where("id = ?", id).First(&sourceComic); sourceComic.Id == 0 {
				continue
			}
			if sourceComic.Retry > 100 {
				continue
			}
			sourceComic.Retry += 1

			recordPick := getPick(sourceComic.Id)
			bot.OnHTML("ol.chapter-page-all", func(e *colly.HTMLElement) {
				e.ForEach(".works-chapter-item", func(sort int, e *colly.HTMLElement) {
					dom := e.DOM.Find("a")
					title, _ := dom.Attr("title")
					url, _ := dom.Attr("href")
					if url == "" {
						return
					}
					aimId, _ := strconv.Atoi(filepath.Base(url))
					if sliceContainsString(recordPick, strconv.Itoa(aimId)) {
						return
					}
					recordPick = append(recordPick, strconv.Itoa(aimId))

					sourceChapter := new(model.SourceChapter)
					sourceChapter.ComicId = sourceComic.Id
					sourceChapter.Source = config.Spe.SourceId
					sourceChapter.Sort = sort
					sourceChapter.Title = title
					sourceChapter.SourceUrl = "https://" + config.Spe.SourceUrl + "/" + strings.TrimLeft(url, "/")
					sourceChapter.SourceChapterId = aimId
					pay := e.DOM.Find("i.ui-icon-pay").Index()
					if pay != -1 {
						sourceChapter.IsFree = 1
					}
					app := e.DOM.Find("span.in-app").Index()
					if app != -1 {
						return
					}
					exists := new(model.SourceChapter).Exists(sourceComic.Id, sourceChapter.SourceUrl)
					if exists == false {
						sourceComic.ChapterPick = sort
						err := orm.Eloquent.Create(&sourceChapter).Error
						if err != nil {
							msg := fmt.Sprintf("chapter数据导入失败 source = %d comic_id = %d chapter_url = %s err = %s",
								config.Spe.SourceId, sourceChapter.ComicId, sourceChapter.SourceUrl, err.Error())
							model.RecordFail(sourceComic.SourceUrl, msg, "漫画章节入库错误", 2)
							rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
						} else {
							rd.RPush(common.SourceChapterTASK, sourceChapter.Id)
							sourceComic.LastChapterUpdateAt = time.Now()
						}
					}
				})
				setPick(sourceComic.Id, recordPick)
			})

			bot.OnHTML("div.works-intro-wr", func(e *colly.HTMLElement) {
				title := e.DOM.Find("h2.works-intro-title").Text()
				state := e.DOM.Find("label.works-intro-status").Text()
				description := e.DOM.Find(".works-intro-short").Text()
				like := e.DOM.Find("#coll_count").Text()
				if state == "已完结" {
					sourceComic.IsFinish = 1
				}
				sourceComic.Description = description
				sourceComic.Like = like
				sourceComic.Title = title
				sourceComic.Region = "国漫"
				sourceComic.SourceData, _ = e.DOM.Html()
				var total int64
				orm.Eloquent.Model(model.SourceChapter{}).Where("comic_id = ?", sourceComic.Id).Count(&total)
				sourceComic.ChapterCount = int(total)
				orm.Eloquent.Save(&sourceComic)
			})

			err := bot.Visit(sourceComic.SourceUrl)
			if err != nil {
				bot = robot.GetColly()
				if i > 5 {
					bot.SetProxy(robot.GetProxy())
				}
				if i == 10 {
					model.RecordFail(sourceComic.SourceUrl, "无法获取漫画详情 :"+sourceComic.SourceUrl, "漫画详情错误", 2)
					rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
				}
			} else {
				break
			}
		}
	}
}

func getPick(comicId int) (recordPick []string) {
	cache := "record:comic:chapters:pick:" + strconv.Itoa(comicId)
	cacheProxy := rd.Get(cache)
	if cacheProxy != "" {
		err := json.Unmarshal([]byte(cacheProxy), &recordPick)
		if err != nil {
			panic(err)
		}
	}
	return recordPick
}

func setPick(comicId int, recordPick []string) {
	cache := "record:comic:chapters:pick:" + strconv.Itoa(comicId)
	jsonData, err := json.Marshal(recordPick)
	if err != nil {
		panic(err)
	}
	rd.Set(cache, string(jsonData), time.Hour*128)
}

func sliceContainsString(slice []string, target string) bool {
	index := sort.SearchStrings(slice, target)
	return index < len(slice) && slice[index] == target
}
