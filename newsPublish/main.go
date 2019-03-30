package main

import (
	_ "newsPublish/routers"
	"github.com/astaxie/beego"
	_"newsPublish/models"
)

func main() {

	beego.AddFuncMap("prePage",ShowPrePage)

	beego.AddFuncMap("nextPage",ShowNextPage)

	beego.Run()

}

//视图函数
//1.在视图中（html）定义一个函数名
//2.在后台定义一个go语言函数
//3.在beego.Run之前把两者关联起来
//
//上一页
func ShowPrePage(pageIndex int )int{
	//pageIndex 是从前端获得的
	if pageIndex<=1{
		return   1
	}
	return pageIndex -1
}

//下一页
func ShowNextPage(pageIndex int ,pageCount float64) int {

	if pageIndex >= int( pageCount){
		return int( pageCount)
	}

	return pageIndex +1
}