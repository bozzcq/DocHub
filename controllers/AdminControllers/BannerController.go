package AdminControllers

import (
	"fmt"

	"time"

	"os"

	"strings"

	"github.com/TruthHun/DocHub/helper"
	"github.com/TruthHun/DocHub/models"
)

//IT文库注册会员管理

type BannerController struct {
	BaseController
}

//横幅列表
func (this *BannerController) Get() {
	var err error
	if this.Data["Banners"], _, err = models.ModelBanner.List(1, 100); err != nil {
		helper.Logger.Error(err.Error())
	}
	this.Data["IsBanner"] = true
	this.TplName = "index.html"
}

//新增横幅
func (this *BannerController) Add() {
	f, h, err := this.GetFile("Picture")
	if err == nil {
		defer f.Close()
		dir := "uploads/" + time.Now().Format("2006-01-02")
		os.MkdirAll(dir, 0777)
		ext := helper.GetSuffix(h.Filename, ".")
		filepath := dir + "/" + helper.MyMD5(fmt.Sprintf("%v-%v", h.Filename, time.Now().Unix())) + "." + ext
		err = this.SaveToFile("Picture", filepath) // 保存位置
		if err == nil {
			if md5str, err := helper.FileMd5(filepath); err == nil {
				save := md5str + "." + ext
				err = models.ModelOss.MoveToOss(filepath, save, true, true)
				if err == nil {
					var banner models.Banner
					this.ParseForm(&banner)
					banner.Picture = save
					banner.TimeCreate = int(time.Now().Unix())
					banner.Status = true
					_, err = models.O.Insert(&banner)
				}
			}
		}
	}
	if err != nil {
		helper.Logger.Error(err.Error())
		this.ResponseJson(0, err.Error())
	}
	this.ResponseJson(1, "横幅添加成功")
}

//删除横幅
func (this *BannerController) Del() {
	var err error
	id := this.GetString("id")
	ids := strings.Split(id, ",")
	if len(ids) > 0 {
		//之所以这么做，是因为如果没有第一个参数，则参数编程了[]string，而不是[]interface{},有疑问可以自己验证试下
		if _, err = models.ModelBanner.Del(ids[0], ids[1:]); err != nil {
			helper.Logger.Error(err.Error())
			this.ResponseJson(0, err.Error())
		}
	}
	this.ResponseJson(1, "删除成功")
}
