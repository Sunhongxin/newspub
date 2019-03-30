package routers

import (
    "github.com/astaxie/beego/context"
    "github.com/astaxie/beego"
    "newsPublish/controllers"

)

func init() {
    //路由过滤器    /article/* : 正则表达式  以article开头的所有路由都需要校验  beego.BeforeExec  执行函数之前,
    // filterFunc :执行函数的函数名
    beego.InsertFilter("/article/*",beego.BeforeExec,filterFunc)

    beego.Router("/", &controllers.MainController{})
    //注册页展示和处理注册业务
    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    //登录页展示和处理登录
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    //展示首页
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    //添加文章
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
    //查看文章详情
    beego.Router("/article/content",&controllers.ArticleController{},"get:Showcontent")
    //编辑文章
    beego.Router("/article/editArticle",&controllers.ArticleController{},"get:ShowEditArticle;post:HandleEditArticle")
    //删除文章
    beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:HandleDelete")
    //添加文章类型
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    //退出登录
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
    //删除文章类型
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
}

//登录状态
func filterFunc(ctx *context.Context){
    //获取session    参数:key值
    userName:= ctx.Input.Session("userName")

    if userName==nil{
        ctx.Redirect(302,"/login")
        return
    }


}