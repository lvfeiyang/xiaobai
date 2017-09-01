package main

import (
	"github.com/lvfeiyang/xiaobai/common/config"
	"github.com/lvfeiyang/xiaobai/common/db"
	"github.com/lvfeiyang/xiaobai/common/flog"
	"github.com/lvfeiyang/xiaobai/message"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
)

var htmlPath string

func main() {
	flog.Init()
	config.Init()
	db.Init()
	htmlPath = config.ConfigVal.HtmlPath

	jsFiles := filepath.Join(htmlPath, "sfk", "js")
	cssFiles := filepath.Join(htmlPath, "sfk", "css")
	fontsFiles := filepath.Join(htmlPath, "sfk", "fonts")
	layDateFiles := filepath.Join(htmlPath, "sfk", "laydate")
	xbjsFiles := filepath.Join(htmlPath, "xiaobai", "html", "js")
	xbcssFiles := filepath.Join(htmlPath, "xiaobai", "html", "css")

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(fontsFiles))))
	http.Handle("/laydate/", http.StripPrefix("/laydate/", http.FileServer(http.Dir(layDateFiles))))
	http.Handle("/xb-js/", http.StripPrefix("/xb-js/", http.FileServer(http.Dir(xbjsFiles))))
	http.Handle("/xb-css/", http.StripPrefix("/xb-css/", http.FileServer(http.Dir(xbcssFiles))))

	http.Handle("/msg/", &message.Message{})

	http.HandleFunc("/xiaobai", xiaobaiHandler)
	flog.LogFile.Fatal(http.ListenAndServe(":80", nil))
}
func xiaobaiHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{
		filepath.Join(htmlPath, "xiaobai", "html", "xb.html"),
		filepath.Join(htmlPath, "xiaobai", "html", "modal", "edit-event.tmpl"),
	}
	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		type oneView struct {
			Id      string
			Time    string
			Address string
			Title   string
			Image   string
			Desc    string
		}
		var view struct {
			EventList []oneView
			WxFlag bool
		}
		if uas, ok := r.Header["User-Agent"]; ok {
			for _, ua := range uas {
				if matched, _ := regexp.MatchString("MicroMessenger.*", ua); matched {
					view.WxFlag = true
				}
			}
		}
		es, err := db.FindAllEvents()
		if err != nil {
			flog.LogFile.Println(err)
		}
		for _, v := range es {
			view.EventList = append(view.EventList, oneView{v.Id.Hex(), v.Time, v.Address, v.Title, message.ImgUrlAddQn(v.Image), v.Desc})
		}
		if err := t.Execute(w, view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
