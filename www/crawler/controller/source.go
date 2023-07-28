package controller

import (
	"comics/controller/kkk"
	"comics/tools/config"
)

type SourceStrategy struct {
	ComicPaw    func()
	ComicUpdate func()
	ChapterPaw  func()
	ImagePaw    func()
}

func SourceOperate(source string) *SourceStrategy {
	switch source {
	case "www.1kkk.com":
		config.Spe.SourceId = 1
		config.Spe.Maxthreads = 1
		return &SourceStrategy{
			ComicPaw:    kkk.ComicPaw,
			ComicUpdate: kkk.ComicUpdate,
			ChapterPaw:  kkk.ChapterPaw,
			ImagePaw:    kkk.ImagePaw,
		}
	}
	return &SourceStrategy{}
}
