package ctemplate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zsxm/scgo/config"
)

//获取所有html模版和被引入的模板
//return ([]string->所有html模版,[]string->被引入的模版)
func Temps(dir, pix, exp string) ([]string, []string) {
	htmlTemps := make([]string, 0, 10)
	includeTemps := make([]string, 0, 10)
	filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, exp) {
			path = path[strings.Index(path, pix):]
			path = strings.Replace(path, "\\", "/", -1)
			htmlTemps = append(htmlTemps, path)
			for _, v := range config.Conf.Template.Include.Files {
				if strings.Index(path, v) != -1 {
					includeTemps = append(includeTemps, path)
				}
			}
		}
		return nil
	})
	return htmlTemps, includeTemps
}
