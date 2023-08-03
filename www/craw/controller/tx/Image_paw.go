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
	"github.com/tebeka/selenium"
	"github.com/ydtg1993/ant"
	"math"
	"regexp"
	"strconv"
	"time"
)

func ImagePaw() {
	A := ant.Get([]int{})
	if A == nil {
		return
	}
	defer A.Free()
	A.Restart("")

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
		if sourceChapter.IsFree == 1 {
			continue
		}

		A.WebDriver.Get(sourceChapter.SourceUrl)
		sourceChapter.Retry += 1
		orm.Eloquent.Save(sourceChapter)
		elem, err := A.WebDriver.FindElement(selenium.ByClassName, "db")
		if err == nil {
			info, err := elem.Text()
			if err == nil {
				if regexp.MustCompile("该漫画不存在或章节已被删除").MatchString(info) {
					continue
				}
			}
		}

		sourceImage := new(model.SourceImage)
		sourceImage.Source = config.Spe.SourceId
		sourceImage.ComicId = sourceChapter.ComicId
		sourceImage.ChapterId = sourceChapter.Id
		sourceImage.Images = model.Images{}
		sourceImage.SourceData = model.Images{}
		browserList(A, sourceImage, sourceChapter)
		if len(sourceImage.SourceData) == 0 {
			msg := fmt.Sprintf("未找到图片资源列表: source = %d comic_id = %d chapter_url = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceUrl)
			A.WebDriver.Refresh()
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
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			} else {
				rd.RPush(common.SourceImageTASK, sourceImage.Id)
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
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片数据更新错误", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			} else {
				rd.RPush(common.SourceImageTASK, record.Id)
			}
		}
	}
}

func browserList(A *ant.Ant, sourceImage *model.SourceImage, sourceChapter *model.SourceChapter) {
	for tryLimit := 0; tryLimit <= 10; tryLimit++ {
		if tryLimit == 6 {
			A.Restart(robot.GetProxy())
			A.WebDriver.Get(sourceChapter.SourceUrl)
		}
		imgContain, err := A.WebDriver.FindElement(selenium.ByClassName, "comic-contain")
		if err != nil {
			if tryLimit == 10 {
				msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl)
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				return
			}
			A.WebDriver.Refresh()
			continue
		}
		var arg []interface{}
		wait := 30
		vh, err := A.WebDriver.ExecuteScript(`return document.getElementById("comicContain").clientHeight`, arg)
		distance := 1200
		distance -= tryLimit * 100
		if err == nil {
			vhi, err := strconv.Atoi(tools.UnknowToString(vh))
			if err == nil {
				wait = int(math.Ceil(float64(vhi) / float64(distance)))
				if wait <= 0 {
					wait = 30
				}
			}
		}
		script := `
		(function(){	
		if (document.getElementById("mainView").scrollTop == 0){
				let f1 = setInterval(()=>{
			 let dom = document.getElementById("mainView")
			 const currentScroll = dom.scrollTop 
			 const clientHeight = dom.clientHeight; 
			 const scrollHeight = dom.scrollHeight; 
			 if (scrollHeight - 10 > currentScroll + clientHeight) {
				 dom.scrollTo({'left':0,'top': currentScroll + %d,behavior: 'smooth'})
			  }else{
				 clearInterval(f1);			
			  }
			},500);
		}else{
			let f2 = setInterval(()=>{
			 let dom = document.getElementById("mainView")
			 const currentScroll = dom.scrollTop 
			 if (currentScroll > 30) {
				 dom.scrollTo({'left':0,'top': currentScroll - %d,behavior: 'smooth'})
			  }else{
				 clearInterval(f2);				
			  }
			},500);
		}
		})();
		;`
		script = fmt.Sprintf(script, distance, distance)
		A.WebDriver.ExecuteScript(script, arg)
		t := time.NewTicker(time.Second * time.Duration(wait+5))
		<-t.C

		imageList, _ := imgContain.FindElements(selenium.ByTagName, "li")
		for _, image := range imageList {
			class, err := image.GetAttribute("class")
			if err == nil && class == "main_ad_top" {
				continue
			}
			img, err := image.FindElement(selenium.ByTagName, "img")
			if err != nil {
				continue
			}
			source, _ := img.GetAttribute("src")
			match, _ := regexp.MatchString("pixel.gif", source)
			if source != "" && match != true {
				sourceImage.SourceData = append(sourceImage.SourceData, source)
			} else if tryLimit < 10 {
				sourceImage.SourceData = model.Images{}
				break
			}
		}
		if len(sourceImage.SourceData) > 0 {
			return
		}
	}
}
