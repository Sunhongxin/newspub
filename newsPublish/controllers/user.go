package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsPublish/models"
)
//继承beego控制器
type UserController struct {
	beego.Controller
}

//展示注册页面
func(this*UserController)ShowRegister(){
	//指定视图页面
	this.TplName = "register.html"
}

//业务处理4小步  1、获取数据 2、校验数据 3、处理数据 （对数据库增删改查）4、返回数据

//处理注册数据
func(this*UserController)HandleRegister(){
//把数据插入到数据库中
	//1、从前端获取数据
	userName :=this.GetString("userName")
	pwd :=this.GetString("password")
	//beego.Info(userName,pwd)  //打印数据为蓝色
	//2、校验数据
	if userName == "" || pwd == "" {
		beego.Error("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}
	//3、操作数据 处理数据 4 小步

	//插入数据
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象  是models包下面的User对象
	var user models.User
	//给插入对象赋值
	user.Name = userName
	user.Pwd = pwd
	//插入
	_,err := o.Insert(&user)
	if err != nil{
		beego.Error("用户注册失败",err)
		this.TplName = "register.html"
		return
	}
	//4、返回数据
	//写一句话给前端
	//this.Ctx.WriteString("用户注册成功")
	//this.TplName = "login.html"
	this.Redirect("/login",302)
}

//展示登录页面
func(this*UserController)ShowLogin(){
	//获取cookie
	userName:=this.Ctx.GetCookie("userName")

	if userName==""{
		this.Data["userName"]=""
		this.Data["checked"]=""
	}else {
		this.Data["userName"]=userName
		this.Data["checked"]="checked"
	}
	//指定视图
	this.TplName = "login.html"
}

//处理登录数据
func(this*UserController)HandleLogin(){
	//获取数据
	userName :=this.GetString("userName")
	pwd :=this.GetString("password")
	//校验数据
	if userName == "" || pwd == ""{
		beego.Error("用户名或者密码不能为空")
		this.TplName = "login.html"
		return
	}
	//处理数据
	//查询
	//获取orm对象
	o := orm.NewOrm()
	//获取查询对象
	var user models.User
	//指定查询条件
	user.Name = userName
	//查询
	err := o.Read(&user,"Name")
	if err != nil{
		beego.Error("用户不存在")
		this.TplName = "login.html"
		return
	}

	//校验密码是否正确
	if user.Pwd != pwd{
		beego.Error("输入的密码错误")
		this.TplName = "login.html"
		return
	}
	//返回数据
	//this.Ctx.WriteString("登录成功")

	//登录成功的情况下,选中复选框把用户名存储到cookie里
	remember:=this.GetString("remember")
	//可以打印一下,看看是什么类型
	if remember=="on"{
		//设置cookie  key值  velue   时间:秒
		this.Ctx.SetCookie("userName",userName,60*60*24)
	}else {
		//取消
		this.Ctx.SetCookie("userName",userName,-1)
	}

	//设置session   key 值   value   在首页展示时获取session
	this.SetSession("userName",userName)

	this.Redirect("/article/index",302)
}

//删除登录状态
func(this *UserController)Logout(){
	//删除session数据
	this.DelSession("userName")
	this.Redirect("/login",302)
}

