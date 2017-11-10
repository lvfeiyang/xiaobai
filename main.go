package main

import (
	"github.com/lvfeiyang/proxy/common"
	"github.com/lvfeiyang/proxy/common/config"
	"github.com/lvfeiyang/proxy/common/db"
	"github.com/lvfeiyang/proxy/common/flog"
	xbDb "github.com/lvfeiyang/xiaobai/common/db"
	"github.com/lvfeiyang/xiaobai/message" //xbMsg
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
)

var htmlPath string
var pjtCfg config.ProjectConfig

func main() {
	flog.Init()
	config.Init()
	db.Init()
	message.Init()
	httpAddr := ":80"
	htmlPath = config.ConfigVal.HtmlPath
	if pjtCfg = config.GetProjectConfig("xiaobai"); "" == pjtCfg.Name {
		flog.LogFile.Fatal("no xiaobai project!")
	}

	if !pjtCfg.Proxy {
		jsFiles := filepath.Join(htmlPath, "sfk", "js")
		cssFiles := filepath.Join(htmlPath, "sfk", "css")
		fontsFiles := filepath.Join(htmlPath, "sfk", "fonts")
		layDateFiles := filepath.Join(htmlPath, "sfk", "laydate")
		http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
		http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(fontsFiles))))
		http.Handle("/laydate/", http.StripPrefix("/laydate/", http.FileServer(http.Dir(layDateFiles))))

		http.Handle("/xiaobai/msg/", &message.LocMessage{})
	} else {
		httpAddr = pjtCfg.Http
	}

	xbjsFiles := filepath.Join(htmlPath, "xiaobai", "html", "js")
	xbcssFiles := filepath.Join(htmlPath, "xiaobai", "html", "css")
	http.Handle("/xiaobai/js/", http.StripPrefix("/xiaobai/js/", http.FileServer(http.Dir(xbjsFiles))))
	http.Handle("/xiaobai/css/", http.StripPrefix("/xiaobai/css/", http.FileServer(http.Dir(xbcssFiles))))

	go common.ListenTcp(pjtCfg.Tcp, message.MhMap)

	http.HandleFunc("/xiaobai", xiaobaiHandler)
	flog.LogFile.Fatal(http.ListenAndServe(httpAddr, nil))
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
			WxFlag    bool
			CanModify bool
		}
		if uas, ok := r.Header["User-Agent"]; ok {
			for _, ua := range uas {
				if matched, _ := regexp.MatchString("MicroMessenger.*", ua); matched {
					view.WxFlag = true
				}
			}
		}
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		user := r.Form.Get("user")
		if "admin" == user {
			view.CanModify = true
		}
		es, err := xbDb.FindAllEvents()
		if err != nil {
			flog.LogFile.Println(err)
		}
		for _, v := range es {
			view.EventList = append(view.EventList, oneView{v.Id.Hex(), v.Time, v.Address, v.Title, common.ImgUrlAddQn(v.Image), v.Desc})
		}
		if err := t.Execute(w, view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
