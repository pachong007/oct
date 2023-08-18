package kk

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"encoding/json"
	"fmt"
	"github.com/shaoerlele/spider"
	"github.com/tebeka/selenium"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ChapterPaw() {
	A := ant.Get([]int{2})
	if A == nil {
		return
	}
	defer A.Free()
	A.Restart("")

	taskLimit := 12
	for limit := 0; limit < taskLimit; limit++ {
		common.StopSignal("章节任务挂起")
		id, err := rd.LPop(common.SourceComicTASK)
		if err != nil || id == "" {
			id, err = rd.LPop(common.SourceComicRenewTASK)
			if err != nil || id == "" {
				return
			}
		}

		var sourceComic model.SourceComic
		if orm.Eloquent.Where("id = ?", id).First(&sourceComic); sourceComic.Id == 0 {
			continue
		}
		if sourceComic.Retry > 100 {
			continue
		}
		sourceComic.Retry += 1

		A.WebDriver.Get(sourceComic.SourceUrl)

		contentBox, err := A.WebDriver.FindElement(selenium.ByClassName, "contentBox")
		if err != nil {
			continue
		}
		contentHTML, err := contentBox.GetAttribute("innerHTML")
		if err != nil {
			continue
		}
		reg := regexp.MustCompile(`class="TopicBox"`)
		if reg.MatchString(contentHTML) == false {
			continue
		}
		reg = regexp.MustCompile(`class="lowershelf"`)
		if reg.MatchString(contentHTML) == true {
			sourceComic.Retry += 100
			continue
		}
		heatDom, err := A.WebDriver.FindElement(selenium.ByClassName, "heat")
		if err == nil {
			popularity, _ := heatDom.Text()
			sourceComic.Popularity = strings.TrimSpace(popularity)
		}

		likeDom, err := A.WebDriver.FindElement(selenium.ByCSSSelector, ".laud .tipTxt")
		if err == nil {
			like, _ := likeDom.GetAttribute("innerText")
			var numbers []rune
			for _, char := range like {
				if unicode.IsDigit(char) {
					numbers = append(numbers, char)
				}
			}
			sourceComic.Like = string(numbers)
		}

		coverDom, err := A.WebDriver.FindElement(selenium.ByClassName, "imgCover")
		if err == nil {
			cover0, _ := coverDom.GetAttribute("data-src")
			cover := strings.TrimSuffix(cover0, "-t.w360.webp.h")
			dir := fmt.Sprintf(config.Spe.DownloadPath+"comic/kk_h_cover/%d/%d", config.Spe.SourceId, sourceComic.SourceId%128)
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
					sourceComic.CoverH = cover
					break
				}
			}
		}

		tabs, err := A.WebDriver.FindElements(selenium.ByClassName, "interval-item")
		if err == nil {
			for page, tab := range tabs {
				tab.Click()
				t := time.NewTicker(time.Second * 1)
				<-t.C
				listElements, _ := A.WebDriver.FindElements(selenium.ByClassName, "title-item")
				chapterList(A, &sourceComic, listElements, page)
			}
		} else {
			listElements, _ := A.WebDriver.FindElements(selenium.ByClassName, "title-item")
			chapterList(A, &sourceComic, listElements, 0)
		}

		detail, err := A.WebDriver.FindElement(selenium.ByClassName, "detailsBox")
		if err == nil {
			sourceComic.Description, _ = detail.Text()
		}
		var total int64
		orm.Eloquent.Model(model.SourceChapter{}).Where("comic_id = ?", sourceComic.Id).Count(&total)
		sourceComic.ChapterCount = int(total)
		orm.Eloquent.Save(&sourceComic)
	}
}

func chapterList(A *ant.Ant, sourceComic *model.SourceComic, listElements []selenium.WebElement, page int) {
	recordPick := getPick(sourceComic.Id)
	for sort, itemElement := range listElements {
		title, _ := itemElement.Text()
		if title == "" {
			continue
		}
		if sliceContainsString(recordPick, title) {
			return
		}
		recordPick = append(recordPick, title)

		sourceChapter := new(model.SourceChapter)
		sourceChapter.Source = 1
		sourceChapter.ComicId = sourceComic.Id
		sourceChapter.Sort = sort + page*50
		sourceChapter.Title = strings.TrimSpace(title)
		html, err := itemElement.GetAttribute("innerHTML")
		if err == nil {
			reg := regexp.MustCompile(`class="lock"`)
			if reg.MatchString(html) {
				continue
			}
		}
		if err := itemElement.Click(); err != nil {
			continue
		}

		allHandles, err := A.WebDriver.WindowHandles()
		if err != nil {
			continue
		}
		if len(allHandles) < 2 {
			continue
		}
		newWindowHandle := allHandles[len(allHandles)-1]
		A.WebDriver.SwitchWindow(newWindowHandle)
		sourceChapter.SourceUrl, _ = A.WebDriver.CurrentURL()
		A.WebDriver.CloseWindow(newWindowHandle)
		A.WebDriver.SwitchWindow(allHandles[0])
		sourceChapter.SourceChapterId = tools.FindStringNumber(sourceChapter.SourceUrl)

		exists := new(model.SourceChapter).Exists(sourceComic.Id, sourceChapter.SourceUrl)
		if exists == false {
			sourceComic.ChapterPick = sort
			err = orm.Eloquent.Create(&sourceChapter).Error
			if err != nil {
				msg := fmt.Sprintf("chapter数据导入失败 source = %d comic_id = %d chapter_url = %s err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl,
					err.Error())
				model.RecordFail(sourceComic.SourceUrl, msg, "漫画章节入库错误", 2)
				rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
			} else {
				rd.RPush(common.SourceChapterTASK, sourceChapter.Id)
				sourceComic.LastChapterUpdateAt = time.Now()
			}
		}
	}
	setPick(sourceComic.Id, recordPick)
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
