package kkk

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
	A := ant.Get([]int{})
	if A == nil {
		return
	}
	defer A.Free()
	if config.Spe.AppDebug == false {
		A.Proxy(robot.GetProxy())
	}

	taskLimit := 9
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
		A.WebDriver.Get(sourceChapter.SourceUrl)
		Mode := 0
		pageDom, err := A.WebDriver.FindElement(selenium.ByClassName, "view-paging")
		if err == nil {
			pageHTML, err := pageDom.GetAttribute("innerHTML")
			if err == nil {
				reg := regexp.MustCompile(`class="chapterpager"`)
				if reg.MatchString(pageHTML) == true {
					//翻页模式
					Mode = 1
				}
			}
		}
		sourceChapter.ViewType = Mode
		sourceChapter.Retry += 1
		orm.Eloquent.Save(sourceChapter)

		ok := saveImgData(sourceChapter)
		if ok == false {
			continue
		}

		dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/%d/%d/%d/%d",
			config.Spe.SourceId, sourceChapter.ComicId%128, sourceChapter.ComicId, sourceChapter.Id)
		sourceImage := new(model.SourceImage)
		if orm.Eloquent.Where("chapter_id = ?", sourceChapter.Id).First(&sourceImage); sourceImage.Id == 0 {
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			continue
		}
		if Mode == 0 {
			browserList(A, sourceImage, sourceChapter)
			downImages(A, sourceChapter, sourceImage, dir)
		} else {
			downImages2(A, sourceChapter, sourceImage, dir)
		}
		err = orm.Eloquent.Save(sourceImage).Error
		if err != nil {
			msg := fmt.Sprintf("图片数据修改失败 source = %d comic_id = %d chapter_id = %d err = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceChapterId,
				err.Error())
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片数据修改", 5)
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			continue
		}
	}
}

func saveImgData(sourceChapter *model.SourceChapter) bool {
	record := new(model.SourceImage)
	exists := record.Exists(sourceChapter.Id)
	if exists == false {
		sourceImage := new(model.SourceImage)
		sourceImage.Source = config.Spe.SourceId
		sourceImage.ComicId = sourceChapter.ComicId
		sourceImage.ChapterId = sourceChapter.Id
		sourceImage.Images = model.Images{}
		sourceImage.SourceData = model.Images{}
		err := orm.Eloquent.Create(&sourceImage).Error
		if err != nil {
			msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceChapterId,
				err.Error())
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误", 2)
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			return false
		}
	} else {
		record.Images = model.Images{}
		record.SourceData = model.Images{}
		record.State = 0
		err := orm.Eloquent.Save(record).Error
		if err != nil {
			msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceChapterId,
				err.Error())
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误.更新", 2)
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			return false
		}
	}

	return true
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func downImages2(A *ant.Ant, sourceChapter *model.SourceChapter, sourceImage *model.SourceImage, dir string) {
	nextButton, err := A.WebDriver.FindElement(selenium.ByXPATH, "//a[@class='block' and contains(@href, 'javascript:ShowNext();')]")
	if err != nil {
		return
	}
outLoop:
	for {
		for try := 0; try <= 3; try++ {
			img, err := A.WebDriver.FindElement(selenium.ByID, "cp_image")
			if try == 3 {
				msg := fmt.Sprintf("未找到图片资源: source = %d comic_id = %d chapter_url = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl)
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 4)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				return
			}
			if err != nil {
				A.WebDriver.Refresh()
				continue
			}
			src, _ := img.GetAttribute("src")
			if containsString(sourceImage.SourceData, src) == false {
				sourceImage.SourceData = append(sourceImage.SourceData, src)
			} else {
				break outLoop
			}
			nextButton.Click()
			break
		}
	}

	if len(sourceImage.SourceData) == 0 {
		msg := fmt.Sprintf("未找到图片资源列表: source = %d comic_id = %d chapter_url = %s",
			config.Spe.SourceId,
			sourceChapter.ComicId,
			sourceChapter.SourceUrl)
		model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
		rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
		sourceImage.Images = []string{}
		sourceImage.State = 0
		return
	}

	files, flag := down(A, sourceImage.SourceData, dir, "jpg", 0)
	if flag == false {
		sourceImage.State = 1
		sourceImage.Images = append(sourceImage.Images, files...)
	} else {
		sourceImage.State = 0
	}
}

func downImages(A *ant.Ant, sourceChapter *model.SourceChapter, sourceImage *model.SourceImage, dir string) {
	imgCount := len(sourceImage.SourceData)
	if imgCount == 0 {
		return
	}
	pageSize := 20
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
			files, flag := down(A, images, dir, "jpg", startIndex+1)
			if flag == false {
				sourceImage.State = 1
				sourceImage.Images = append(sourceImage.Images, files...)
			} else {
				sourceImage.State = 0
			}
			browserList(A, sourceImage, sourceChapter)
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
	files, flag := down(A, sourceImage.SourceData, dir, "jpg", 0)
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

func browserList(A *ant.Ant, sourceImage *model.SourceImage, sourceChapter *model.SourceChapter) {
	for tryLimit := 0; tryLimit <= 9; tryLimit++ {
		imgList, err := A.WebDriver.FindElements(selenium.ByClassName, "load-src")
		if err != nil {
			if tryLimit > 5 {
				if tryLimit == 9 {
					msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %s err = %s",
						config.Spe.SourceId,
						sourceChapter.ComicId,
						sourceChapter.SourceUrl,
						err.Error())
					model.RecordFail(sourceChapter.SourceUrl, msg, "图片列表未找到", 3)
					return
				}
				A.Proxy(robot.GetProxy())
			}
			continue
		}
		sourceImage.SourceData = []string{}
		for _, img := range imgList {
			img, err := img.GetAttribute("data-src")
			if err == nil {
				sourceImage.SourceData = append(sourceImage.SourceData, img)
			}
		}
		if len(sourceImage.SourceData) > 0 {
			return
		}
		A.WebDriver.Refresh()
		A.WebDriver.Get(sourceChapter.SourceUrl)
	}
}
