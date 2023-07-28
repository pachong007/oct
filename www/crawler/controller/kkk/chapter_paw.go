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
	"regexp"
	"strings"
	"time"
)

type DomChapter struct {
	Title string
	Url   string
	Id    int
	Cover string
}

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
			if sourceComic.Retry > 30 {
				continue
			}
			sourceComic.Retry += 1

			bot.OnHTML("#tempc", func(e *colly.HTMLElement) {
				order := e.DOM.Find(".order").Text()

				var Chapters []*DomChapter
				e.ForEach("#chapterlistload li", func(sort int, e *colly.HTMLElement) {
					dom := e.DOM.Find("a")
					title := strings.TrimSpace(dom.Text())
					url, _ := dom.Attr("href")
					if e.DOM.HasClass("detail-lock") {
						return
					}

					Ch := new(DomChapter)
					Ch.Title = title
					Ch.Url = url
					re := regexp.MustCompile(`/.*-(\d+)/`)
					match := re.FindStringSubmatch(url)
					if len(match) > 1 {
						numberStr := match[1]
						var number int
						_, err := fmt.Sscanf(numberStr, "%d", &number)
						if err == nil {
							Ch.Id = number
						}
					}

					cover, ok := e.DOM.Find("img.img").Attr("src")
					if ok {
						var cookies map[string]string
						dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/cover/%d/%d", config.Spe.SourceId, sourceComic.Id%128)
						for tryLimit := 0; tryLimit <= 7; tryLimit++ {
							proxy := ""
							if tryLimit > 5 {
								proxy = robot.GetProxy()
							}
							downCover := common.DownFile(cover, dir, tools.RandStr(9)+".jpg", proxy, cookies)
							if downCover != "" {
								Ch.Cover = downCover
								break
							}
						}
					}
					if url == "" || Ch.Id == 0 {
						return
					}
					Chapters = append(Chapters, Ch)
				})
				if order == "倒序" {
					reverseSlice(Chapters)
				}

				for sort, chapter := range Chapters {
					if sort < sourceComic.ChapterPick {
						continue
					}
					sourceChapter := new(model.SourceChapter)
					sourceChapter.ComicId = sourceComic.Id
					sourceChapter.Source = config.Spe.SourceId
					sourceChapter.Sort = sort
					sourceChapter.Title = chapter.Title
					sourceChapter.SourceUrl = "https://" + config.Spe.SourceUrl + "/" + strings.TrimLeft(chapter.Url, "/")
					sourceChapter.SourceChapterId = chapter.Id
					sourceChapter.Cover = chapter.Cover
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
				}
			})

			bot.OnHTML("div.banner_detail_form", func(e *colly.HTMLElement) {
				sourceComic.Cover, _ = e.DOM.Find(".cover>img").Attr("src")
				var cookies map[string]string
				dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/%d/%d", config.Spe.SourceId, sourceComic.Id%128)
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

				author := e.DOM.Find("p.subtitle>a").Text()
				state := e.DOM.Find("p.tip>span>span").Text()
				description := e.DOM.Find("p.content").Text()

				if state == "连载中" {
					sourceComic.IsFinish = 0
				} else {
					sourceComic.IsFinish = 1
				}
				sourceComic.Author = author
				sourceComic.Description = description
				sourceComic.SourceData, _ = e.DOM.Html()
				var total int64
				orm.Eloquent.Model(model.SourceChapter{}).Where("comic_id = ?", sourceComic.Id).Count(&total)
				sourceComic.ChapterCount = int(total)
				orm.Eloquent.Save(&sourceComic)
			})

			err := bot.Visit(sourceComic.SourceUrl)
			if err != nil {
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

func reverseSlice(s []*DomChapter) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
