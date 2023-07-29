package main

import (
	"comics/common"
	"comics/controller"
	"comics/robot"
	"comics/robot/ant"
	"comics/tools/config"
	"comics/tools/database"
	"comics/tools/log"
	"comics/tools/rd"
	"os"
	"time"
)

var Source *controller.SourceStrategy

func main() {
	Setup()
	t := time.NewTicker(time.Second * 10)
	<-t.C
	go controller.TaskComic(Source)
	go controller.TaskChapter(Source)
	go controller.TaskChapterUpdate()
	go controller.TaskReComic(Source)
	controller.TaskImage(Source)
}

func Setup() {
	err := config.Spe.SetUp()
	if err != nil {
		panic(err)
	}

	url := os.Getenv("SOURCE_URL")
	if url != "" {
		config.Spe.SourceUrl = url
	}
	Source = controller.SourceOperate(config.Spe.SourceUrl)
	config.Spe.RedisDb = config.Spe.SourceId

	mylog := new(log.LogsManage)
	err = mylog.SetUp()
	if err != nil {
		panic(err)
	}

	db := new(database.MysqlManage)
	err = db.Setup()
	if err != nil {
		panic(err)
	}

	redisManage := new(rd.RedisManage)
	err = redisManage.SetUp()
	if err != nil {
		panic(err)
	}
	rd.Delete(common.TaskStepRecord)
	rd.Delete(common.StopRobotSignal)

	go ant.Build(config.Spe.Maxthreads, robot.GetSeleniumArgs())
}
