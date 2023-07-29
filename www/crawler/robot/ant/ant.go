package ant

import (
	"archive/zip"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var SeleniumPath string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("获取文件路径失败")
		return
	}
	absPath, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		fmt.Println("获取文件绝对路径失败：", err)
		return
	}
	file := ""

	Sys := runtime.GOOS
	if Sys == "windows" {
		file = absPath + "\\driver-win.zip"
		_, err := os.Stat(absPath + "/" + "chromedriver.exe")
		if err == nil {
			SeleniumPath = absPath + "\\" + "chromedriver.exe"
			return
		}
	} else {
		file = absPath + "/driver.zip"
		_, err := os.Stat(absPath + "/" + "chromedriver")
		if err == nil {
			SeleniumPath = absPath + "/" + "chromedriver"
			return
		}
	}

	r, err := zip.OpenReader(file)
	if err != nil {
		fmt.Println("Failed to open package:", err)
		return
	}
	defer r.Close()

	exeFileMode := os.FileMode(0777)
	for _, f := range r.File {
		if strings.Contains(f.Name, "chromedriver") {
			exeReader, err := f.Open()
			if err != nil {
				fmt.Println("Failed to open executable:", err)
				return
			}
			defer exeReader.Close()
			exeFileWriter, err := os.OpenFile(absPath+"/"+f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, exeFileMode)
			if err != nil {
				fmt.Println("Failed to open executable file:", err)
				return
			}
			defer exeFileWriter.Close()
			_, err = io.Copy(exeFileWriter, exeReader)
			SeleniumPath = absPath + "/" + f.Name
			break
		}
	}
}

var Swarm []*Ant

type Ant struct {
	Service   *selenium.Service
	WebDriver selenium.WebDriver
	Port      int
	Headers   map[string]string
	Lock      sync.Mutex
	IsLocked  bool
}

func Build(num int, headers map[string]string) {
	for i := 0; i < num; i++ {
		ant := &Ant{
			Headers:  mergeHeader(headers),
			Lock:     sync.Mutex{},
			IsLocked: false,
		}
		ant.setPort()
		ant.prepare()
		Swarm = append(Swarm, ant)
	}
}

func Get(pick []int) (A *Ant) {
	for index, ant := range Swarm {
		if pickIndex(index, pick) == false {
			continue
		}
		if ant.IsLocked == true {
			continue
		}
		ant.Lock.Lock()
		ant.IsLocked = true
		A = ant
		break
	}
	return A
}

func pickIndex(index int, pick []int) bool {
	if len(pick) == 0 {
		return true
	}
	result := false
	if pick[0] < 0 {
		result = true
	}
	for _, p := range pick {
		if p == 0 {
			continue
		}
		if result == false && p > 0 {
			p--
			if index == p {
				return true
			}
			continue
		}
		if result == true && p < 0 {
			p++
			if index == int(math.Abs(float64(p))) {
				return false
			}
		}
	}
	return result
}

func (ant *Ant) Free() {
	ant.IsLocked = false
	ant.Lock.Unlock()
}

func (ant *Ant) Proxy(proxy string) {
	ant.WebDriver.Close()
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	ant.Headers["--proxy-server"] = proxy
	var header = []string{}
	for k, v := range ant.Headers {
		var arg string
		if v == "" {
			arg = k
		} else {
			arg = k + "=" + v
		}
		header = append(header, arg)
	}

	caps.AddChrome(chrome.Capabilities{
		Prefs: map[string]interface{}{
			"profile.managed_default_content_settings.images":             2,
			"profile.default_content_setting_values.notifications":        1,
			"profile.content_settings.exceptions.notifications.*.setting": 1,
		},
		Args: header,
	})
	wb, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", ant.Port))
	if err != nil {
		ant.prepare()
		return
	}
	ant.WebDriver = wb
}

func mergeHeader(args map[string]string) map[string]string {
	args0 := map[string]string{
		"--headless":                  "",
		"--no-sandbox":                "",
		"--ignore-certificate-errors": "",
		"--ignore-ssl-errors":         "",
		"--user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}
	for k, v := range args {
		if v == "-" {
			delete(args0, k)
			continue
		}
		args0[k] = v
	}
	return args0
}

func (ant *Ant) setPort() {
	ant.Port += rand.Intn(99) + 1
	if ant.Port > 55000 {
		ant.Port = ant.Port - 35000 - rand.Intn(999)
	}
}

func (ant *Ant) prepare() {
	for tryLimit := 0; tryLimit <= 999; tryLimit++ {
		if ant.WebDriver != nil {
			ant.WebDriver.Close()
		}
		if ant.Service != nil {
			ant.Service.Stop()
		}
		opts := []selenium.ServiceOption{
			selenium.ChromeDriver(SeleniumPath),
		}
		service, err := selenium.NewChromeDriverService(SeleniumPath, ant.Port, opts...)
		if nil != err {
			ant.setPort()
			if tryLimit == 999 {
				panic(err.Error())
			}
			continue
		}
		ant.Service = service

		caps := selenium.Capabilities{
			"browserName": "chrome",
		}

		var header = []string{}
		for k, v := range ant.Headers {
			var arg string
			if v == "" {
				arg = k
			} else {
				arg = k + "=" + v
			}
			header = append(header, arg)
		}

		caps.AddChrome(chrome.Capabilities{
			Prefs: map[string]interface{}{
				"profile.managed_default_content_settings.images":             2,
				"profile.default_content_setting_values.notifications":        1,
				"profile.content_settings.exceptions.notifications.*.setting": 1,
			},
			Args: header,
		})
		wb, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", ant.Port))
		if err != nil {
			if tryLimit == 999 {
				panic(err.Error())
			}
			continue
		}
		wb.SetImplicitWaitTimeout(time.Second * 60)
		wb.SetPageLoadTimeout(time.Second * 60)
		wb.ResizeWindow("", 1400, 1200)
		ant.WebDriver = wb
		break
	}
}
