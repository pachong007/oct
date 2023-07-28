package controller

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func TaskComic(source *SourceStrategy) {
	t := time.NewTicker(time.Hour * 12)
	defer t.Stop()
	rd.RPush(common.TaskStepRecord, fmt.Sprintf("漫画-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
	source.ComicPaw()
	source.ComicPaw()
	for {
		<-t.C
		rd.Delete(common.TaskStepRecord)
		rd.RPush(common.TaskStepRecord, fmt.Sprintf("漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		source.ComicUpdate()
	}
}

func TaskReComic(source *SourceStrategy) {
	t := time.NewTicker(time.Hour * 360)
	defer t.Stop()

	for {
		<-t.C
		source.ComicPaw()
	}
}

func TaskChapter(source *SourceStrategy) {
	t := time.NewTicker(time.Minute * 5)
	defer t.Stop()
	threads := 5
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(threads)
		rd.RPush(common.TaskStepRecord, fmt.Sprintf("章节-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < threads; i++ {
			go func() {
				source.ChapterPaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TaskChapterUpdate() {
	t := time.NewTicker(time.Hour * 6)
	defer t.Stop()
	for {
		<-t.C
		rd.RPush(common.TaskStepRecord, fmt.Sprintf("连载漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		p := rd.Get(common.SourceComicRenewPick)
		var page int
		if p == "" {
			page = 0
		} else {
			page, _ = strconv.Atoi(p)
		}
		limit := 1000

		var sourceComics []model.SourceComic
		orm.Eloquent.Offset(page*limit).Limit(limit).Where("source = ? and is_finish = 0", config.Spe.SourceId).Find(&sourceComics)
		if len(sourceComics) <= 0 {
			rd.Set(common.SourceComicRenewPick, "0", time.Hour*9999)
			continue
		}
		page = page + 1
		for _, sourceComic := range sourceComics {
			rd.RPush(common.SourceComicRenewTASK, sourceComic.Id)
		}
		rd.Set(common.SourceComicRenewPick, strconv.Itoa(page), time.Hour*9999)
	}
}

func TaskImage(source *SourceStrategy) {
	for {
		timestamp0 := time.Now()
		threads := 10
		wg := sync.WaitGroup{}
		wg.Add(threads)
		for i := 0; i < threads; i++ {
			t0 := time.NewTicker(time.Second * 30)
			<-t0.C
			go func() {
				rd.Set(common.SourceImageCapture,
					fmt.Sprintf("图片链接抓取 %s", time.Now().String()),
					time.Hour*1)
				source.ImagePaw()
				rd.Set(common.SourceImageDownload,
					fmt.Sprintf("图片下载 %s", time.Now().String()),
					time.Hour*1)
				if config.Spe.SourceId == 2 {
					common.DownImage("jpg", 10)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		if time.Now().Sub(timestamp0) < 5*time.Second {
			t := time.NewTicker(time.Minute * 15)
			<-t.C
		}
	}
}
