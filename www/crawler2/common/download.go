package common

import (
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

func DownImage(ext string, taskLimit int) {
	for limit := 0; limit < taskLimit; limit++ {
		StopSignal("图片下载任务挂起")
		id, err := rd.LPop(SourceImageTASK)
		if err != nil || id == "" {
			return
		}
		sourceImage := new(model.SourceImage)
		if orm.Eloquent.Where("id = ?", id).First(&sourceImage); sourceImage.Id == 0 {
			continue
		}
		sourceChapter := new(model.SourceChapter)
		orm.Eloquent.Where("id = ?", sourceImage.ChapterId).First(&sourceChapter)
		if sourceChapter.IsFree == 1 {
			continue
		}
		dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/%d/%d/%d/%d",
			config.Spe.SourceId, sourceChapter.ComicId%128, sourceChapter.ComicId, sourceImage.ChapterId)
		down(sourceImage, dir, ext)
		if len(sourceImage.Images) == 0 {
			logs.Warning(fmt.Sprintf("图片下载本地失败 source_id = %d comic_id = %d chapter_id = %d",
				config.Spe.SourceId, sourceChapter.ComicId, sourceChapter.Id))
			rd.RPush(SourceChapterRetryTask, sourceChapter.Id)
			continue
		}
		orm.Eloquent.Save(&sourceImage)
	}
}

func down(sourceImage *model.SourceImage, dir, ext string) {
	for key, img := range sourceImage.SourceData {
		sourceImage.State = 0
		proxy := ""
		for tryLimit := 0; tryLimit <= 5; tryLimit++ {
			file := DownFile(img, dir, fmt.Sprintf("%d.%s", key, ext), proxy, map[string]string{})
			if file != "" {
				sourceImage.State = 1
				sourceImage.Images = append(sourceImage.Images, file)
				break
			} else if tryLimit > 3 {
				proxy = robot.GetProxy()
			}
		}
	}
	if sourceImage.State == 0 {
		sourceImage.Images = model.Images{}
	}
}
