package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	IP          string
	Port        uint64
	NameSpaceID string
	DataID      string
	Group       string
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Print("获取当前程序执行目录失败")
		return
	}

	file, err := ini.Load(fmt.Sprintf("%s/%s", dir, "conf/config.ini"))
	if err != nil {
		fmt.Println("配置文件读取有误,请检查配置文件.", err)
		return
	}

	LoadNacos(file)
}

func LoadNacos(file *ini.File) {
	s := file.Section("nacos")
	IP = s.Key("IP").String()
	Port = s.Key("Port").MustUint64(8848)
	NameSpaceID = s.Key("NameSpaceID").String()
	DataID = s.Key("DataID").String()
	Group = s.Key("Group").String()
}
