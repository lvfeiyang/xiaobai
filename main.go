package main

import (
	"github.com/lvfeiyang/xiaobai/common/flog"
	"github.com/lvfeiyang/xiaobai/common/config"
	"github.com/lvfeiyang/xiaobai/common/db"
	"github.com/lvfeiyang/xiaobai/message"
	"net/http"
	"path/filepath"
	"html/template"
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
	xbjsFiles := filepath.Join(htmlPath, "xiaobai", "html", "js")
	xbcssFiles := filepath.Join(htmlPath, "xiaobai", "html", "css")

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(fontsFiles))))
	http.Handle("/xb-js/", http.StripPrefix("/xb-js/", http.FileServer(http.Dir(xbjsFiles))))
	http.Handle("/xb-css/", http.StripPrefix("/xb-css/", http.FileServer(http.Dir(xbcssFiles))))

	http.Handle("/msg/", &message.Message{})

	http.HandleFunc("/xiaobai", xiaobaiHandler)
	flog.LogFile.Fatal(http.ListenAndServe(":7070", nil))
}
func xiaobaiHandler(w http.ResponseWriter, r *http.Request)  {
	paths := []string{
		filepath.Join(htmlPath, "xiaobai", "html", "xb.html"),
		// filepath.Join(htmlPath, "xiaobai", "html", "modal", "edit-event.tmpl"),
	}
	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		if err := t.Execute(w, nil); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
