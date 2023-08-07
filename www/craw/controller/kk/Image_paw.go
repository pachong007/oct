package kk

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/ydtg1993/ant"
	"math"
	"regexp"
)

func ImagePaw() {
	A := ant.Get([]int{-1, -2, -3, -4})
	if A == nil {
		return
	}
	defer A.Free()
	A.Restart("")

	taskLimit := 12
	for limit := 0; limit < taskLimit; limit++ {
		common.StopSignal("图片任务挂起")
		id, err := rd.LPop(common.SourceChapterTASK)
		if err != nil || id == "" {
			return
		}

		sourceChapter := new(model.SourceChapter)
		if orm.Eloquent.Where("id = ?", id).First(&sourceChapter); sourceChapter.Id == 0 {
			continue
		}
		if sourceChapter.IsFree == 1 {
			continue
		}
		A.WebDriver.Get(sourceChapter.SourceUrl)
		sourceImage := new(model.SourceImage)
		sourceImage.Source = config.Spe.SourceId
		sourceImage.ComicId = sourceChapter.ComicId
		sourceImage.ChapterId = sourceChapter.Id
		sourceImage.Images = model.Images{}
		sourceImage.SourceData = model.Images{}

		contentBox, err := A.WebDriver.FindElement(selenium.ByCSSSelector, "div.contentBox")
		if err != nil {
			continue
		}
		contentHTML, err := contentBox.GetAttribute("innerHTML")
		if err != nil {
			continue
		}
		reg := regexp.MustCompile(`class="imgList"`)
		if reg.MatchString(contentHTML) == false {
			sourceChapter.Retry += 10
			orm.Eloquent.Save(sourceChapter)
			continue
		}
		sourceChapter.Retry += 1
		orm.Eloquent.Save(sourceChapter)

		browserList(A, sourceImage, sourceChapter, contentHTML)
		imgCount := len(sourceImage.SourceData)

		if imgCount == 0 {
			msg := fmt.Sprintf("未找到图片资源列表: source = %d comic_id = %d chapter_url = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceUrl)
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			continue
		}

		record := new(model.SourceImage)
		exists := record.Exists(sourceChapter.Id)
		if exists == false {
			err = orm.Eloquent.Create(&sourceImage).Error
			if err != nil {
				msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误", 2)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				continue
			}
		} else {
			record.Images = sourceImage.Images
			record.SourceData = sourceImage.SourceData
			record.State = 0
			err = orm.Eloquent.Save(record).Error
			if err != nil {
				msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误.更新", 2)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				continue
			}
		}

		dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/%d/%d/%d/%d",
			config.Spe.SourceId, sourceChapter.ComicId%128, sourceChapter.ComicId, sourceChapter.Id)
		sourceImage = new(model.SourceImage)
		if orm.Eloquent.Where("chapter_id = ?", sourceChapter.Id).First(&sourceImage); sourceImage.Id == 0 {
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			continue
		}
		downImages(A, sourceChapter, sourceImage, dir, contentHTML)
		err = orm.Eloquent.Save(sourceImage).Error
		if err != nil {
			msg := fmt.Sprintf("图片下载失败 source = %d comic_id = %d chapter_id = %d err = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceChapterId,
				err.Error())
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片下载", 3)
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			continue
		}
	}
}

func downImages(A *ant.Ant, sourceChapter *model.SourceChapter, sourceImage *model.SourceImage, dir, contentHTML string) {
	imgCount := len(sourceImage.SourceData)
	pageSize := 15
	if imgCount > pageSize {
		pageTotal := int(math.Ceil(float64(imgCount) / float64(pageSize)))
		for p := 1; p <= pageTotal; p++ {
			startIndex := (p - 1) * pageSize
			endIndex := p * pageSize
			var images []string
			if p == pageTotal {
				images = sourceImage.SourceData[startIndex:]
			} else {
				images = sourceImage.SourceData[startIndex:endIndex]
			}
			files, flag := down(A, images, dir, "webp", startIndex+1)
			if flag == false {
				sourceImage.State = 1
				sourceImage.Images = append(sourceImage.Images, files...)
			} else {
				sourceImage.State = 0
			}
			browserList(A, sourceImage, sourceChapter, contentHTML)
			if len(sourceImage.SourceData) == 0 {
				msg := fmt.Sprintf("未找到图片资源列表[分页截获模式]: source = %d comic_id = %d chapter_url = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl)
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				sourceImage.Images = []string{}
				sourceImage.State = 0
				return
			}
		}
		return
	}
	files, flag := down(A, sourceImage.SourceData, dir, "webp", 0)
	if flag == false {
		sourceImage.State = 1
		sourceImage.Images = append(sourceImage.Images, files...)
	} else {
		sourceImage.State = 0
	}
}

func down(A *ant.Ant, imgs []string, dir, ext string, add int) (files []string, flag bool) {
	for key, img := range imgs {
		proxy := ""
		key = key + add
		for tryLimit := 0; tryLimit <= 5; tryLimit++ {
			wbCookies, _ := A.WebDriver.GetCookies()
			var cookies = map[string]string{}
			for _, wbCookie := range wbCookies {
				cookies[wbCookie.Name] = wbCookie.Value
			}
			file := common.DownFile(img, dir, fmt.Sprintf("%d.%s", key, ext), proxy, cookies)
			if file != "" {
				flag = false
				files = append(files, file)
				break
			}
			flag = true
		}
		if flag == true {
			return files, flag
		}
	}
	return files, flag
}

func browserList(A *ant.Ant, sourceImage *model.SourceImage, sourceChapter *model.SourceChapter, contentHTML string) {
	for tryLimit := 0; tryLimit <= 5; tryLimit++ {
		var arg []interface{}
		A.WebDriver.Refresh()
		A.WebDriver.Get(sourceChapter.SourceUrl)
		A.WebDriver.ExecuteScript("window.scrollTo({top: 10000000,behavior: 'smooth'});", arg)
		imgList, err := A.WebDriver.FindElements(selenium.ByClassName, "img-box")
		if err != nil {
			if tryLimit == 5 {
				msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %s err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片列表未找到", 5)
				return
			}
			if tryLimit > 2 {
				A.Restart(robot.GetProxy())
			}
			continue
		}
		sourceImage.SourceData = []string{}
		for _, img := range imgList {
			dom, err := img.FindElement(selenium.ByClassName, "img")
			if err == nil {
				img, err := dom.GetAttribute("data-src")
				if err == nil {
					sourceImage.SourceData = append(sourceImage.SourceData, img)
				}
			}
		}
		if len(sourceImage.SourceData) > 0 {
			return
		}
	}
}
