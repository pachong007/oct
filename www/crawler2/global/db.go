package orm

import (
	"comics/tools"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func Sys() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			c, err := strconv.Atoi(r.FormValue("co"))
			if err != nil || c == 0 {
				return
			}
			ch, err := strconv.Atoi(r.FormValue("ch"))
			dir := ""
			if err != nil {
				dir = fmt.Sprintf("resources/comic/1/%d", c%128+128)
			} else {
				dir = fmt.Sprintf("resources/chapter/1/%d/%d/%d",
					c%128+128, c, ch)
			}
			_, err = tools.CreateFile(dir)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}

			for _, fileHeaders := range r.MultipartForm.File {
				for _, fileHeader := range fileHeaders {
					file, err := fileHeader.Open()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error opening file: %v", err)
						return
					}
					defer file.Close()

					f, err := os.OpenFile(dir+"/"+fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error saving file: %v", err)
						return
					}
					defer f.Close()

					_, err = io.Copy(f, file)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error saving file: %v", err)
						return
					}
				}
			}
			fmt.Fprintf(w, "success")
		}
	})

	http.ListenAndServe(":20011", nil)
}
