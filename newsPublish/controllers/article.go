package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"newsPublish/models"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}


//展示首页内容
func(this*ArticleController)ShowIndex(){
	//获取session
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/article/login",302)
		return
	}


	//获取所有文章数据
	//1.获取orm对象
	o := orm.NewOrm()
	//获取所有文章    select * from article;  queryseter
	//2.指定表Article
	// 从 Article表中获取所有 文章数据 返回一个对象
	qs :=o.QueryTable("Article")

	//声明 所有文章数据存储的变量
	var articles []models.Article
	//qs.All(&articles)//获取所有
	//分页实现

	pageSize:=2//每页记录数


	//处理首页和末页
	pageIndex,err:=this.GetInt("pageIndex")
	if err!=nil{
		pageIndex=1
	}

	//获取数据库部分数据  qs.Limit(参1,参2)参1:获取几条数据    参2:从哪获取
	start:=pageSize*(pageIndex-1)
	//获取所有文章类型数据给前端
	//定义一个变量存储这些文章类型
	var   count int64
	var articleTypes  []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	//传递给前端
	this.Data["articleTypes"]=articleTypes
	//下拉框改变的时候获取相应类型的文章
	typeName:=this.GetString("select")
	//下拉框有没有数据
	if typeName=="" {
		//获取总记录数和总页数
		count,_=qs.RelatedSel("ArticleType").Count()//总记录数
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else {
		//获取到类型数据,根据这个数据获得相应的文章
		//多表查询是惰性查询  RelatedSel左联查询
		//Filter("ArticleType__TypeName",typeName)  相当于where
		//获取总记录数和总页数
		count,_=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()//总记录数
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

	}
	//总页数
	pageCount:=math.Ceil(float64(float64(count)/float64(pageSize)))//页数  天花板函数
	this.Data["pageCount"]=pageCount
	this.Data["Count"]=count

	this.Data["pageIndex"]=pageIndex
	this.Data["typeName"]=typeName
	//传递数据给前端
	this.Data["articles"] = articles

	//展示layout
	this.Layout="layout.html"

	//指定视图 展示首页
	this.TplName = "index.html"
}

//展示添加文章页面
func(this*ArticleController)ShowAdd(){
	//获取所有类型数据,传递给前端
	//获取orm对象
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	//根据文章类型获取所有文章类型
	o.QueryTable("articleType").All(&articleTypes)
	//把所有的文章类型传递给前端
	this.Data["articleTypes"]=articleTypes

	//展示layout
	this.Layout="layout.html"

	//指定视图
	this.TplName = "add.html"
}
//
//处理添加文章数据
func(this*ArticleController)HandleAdd(){
	//从前端 获取数据
	articleName :=this.GetString("articleName")//文章名字
	content :=this.GetString("content")//文章内容
	beego.Info(articleName,content)
	//func (c *Controller) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	//file字节流；FileHeader文件头是一个结构体 Filename Header Size content 等
	file,head,err :=this.GetFile("uploadname")//图片文件

	//获取数据
	if articleName == "" || content == "" || err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	//需要图片判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return
	}
	//需要校验图片格式   截取文件名  path.Ext  获取文件名.  点后面的内容  如：.jpg  .png  .jpeg
	ext :=path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return
	}

	//防止图片重名
	//beego.Info("打印当前时间CST格式的时间：time.now = ",time.Now().Format("2006-01-02 15:04:05"))
	fileName := time.Now().Format("20060102150405")
	//操作数据  SaveToFile：beego存储文件数据的方法    fileName：格式化时间  ext：后缀

	this.SaveToFile("uploadname","./static/img/"+fileName+ext)
	//beego.Info(head.Filename)//检验错误位置
	//把数据插入到数据库
	//获取orm对象
	o := orm.NewOrm()
	//获取插入对象
	var article models.Article
	//给插入对象赋值
	article.Title = articleName
	article.Content = content
	article.Img = "/static/img/"+fileName+ext //文件路径


	//获取文章类型数据
	//从前端拿去文章类型数据
	typeName:=this.GetString("select")
	//获取类型数据的对象
	var articleType models.ArticleType
	//赋值
	articleType.TypeName=typeName
	//读取文章类型的类型名
	o.Read(&articleType,"TypeName")
	//给文章的文章类型赋值
	article.ArticleType=&articleType

	//插入
	_,err=o.Insert(&article)
	if err!=nil{
		beego.Error("插入失败",err)
		return
	}
	//返回数据
	this.Redirect("/article/index",302)
}

//展示文章详情
func (this*ArticleController)Showcontent(){
 	//获取数据
 	id,err:=this.GetInt("id")

 	//校验数据
 	if err!=nil{
 		beego.Error("获取文章详情失败",err)
 		this.TplName="index.html"
		return
	}

 	//处理数据 ,查询数据,获取文章数据
 	//获取orm对象
 	o:=orm.NewOrm()

 	//获取查询对象
 	var alticle models.Article

 	//指定查询条件
 	alticle.Id2=id
 	//执行查询查询操作
 	o.Read(&alticle)
 	//多对多查询的两种方式
 	//o.LoadRelated(&alticle,"Users")//多对多查询在article表里查询User字段   有重复的
 	var users []models.User
 	//user 表    Articles:user表多对多字段名  __Article: 类型名/表明 __Id2  :安主键查   alticle.Id2:要比较的条件
 	// Distinct()  过滤去重
 	o.QueryTable("user").Filter("Articles__Article__Id2",alticle.Id2).Distinct().All(&users)//高级查询
 	this.Data["users"]=users


 	//阅读次数+1  更新操作
 	alticle.ReadCount+=1
 	//更新
 	o.Update(&alticle)

 	//返回数据给前端
	this.Data["article"]=alticle
	//添加浏览记录   多对多插入   Users:表中的字段
	m2m:=o.QueryM2M(&alticle,"Users")
	//向这个表里插入对象指针,获取对象

	var user models.User
	//获取session里的用户名
	userName:=this.GetSession("userName")
	//userName.(string)  :GetSession  返回接口类型  需要断言
	//user.Name   user表  Name字段
    user.Name=userName.(string)//从session中获得用户名
    o.Read(&user,"Name")//获取当前用户

    //插入对象
    m2m.Add(user)


	//展示layout
	this.Layout="layout.html"
	//this.LayoutSections=make(map[string]string)	//加载js
	//this.LayoutSections["JsFile"]="index.js"


	//指定视图
	this.TplName="content.html"
}

//展示文章编辑页
func (this*ArticleController)ShowEditArticle(){
	//填充的文章原来的数据
	// 获取数据
	id,err:=this.GetInt("id")
	beego.Info(id,err)
	//校验数据
	if err!=nil{
		beego.Error("获取文章数据失败",err)
		this.TplName="index.html"
		return
	}
	//处理数据
	//查询
	//1 获取orm对象
	o:=orm.NewOrm()
	//2 指定查询对象
	var article models.Article
	//3 指定查询条件
	article.Id2=id
	//4 查询 Read
	o.Read(&article)

	//返回数据
	this.Data["article"]=article

	//展示layout
	this.Layout="layout.html"

	//指定视图
	this.TplName="update.html"
}

//####封装一个文件校验的函数
func UploadFunc(this *ArticleController,FileName string)(string){
	file,head,err :=this.GetFile(FileName)//图片文件

	//获取数据
	if err != nil{
		beego.Error("获取用户添加数据失败",err)
		this.TplName = "add.html"
		return "" //错误返回空
	}
	defer file.Close()

	//需要图片判断大小
	if head.Size > 5000000{
		beego.Error("图片太大，我不收")
		this.TplName = "add.html"
		return "" //错误返回空
	}
	//需要校验图片格式   截取文件名  path.Ext  获取文件名.  点后面的内容  如：.jpg  .png  .jpeg
	ext :=path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		beego.Error("文件格式不正确")
		this.TplName = "add.html"
		return "" //错误返回空
	}

	//防止图片重名
	//beego.Info("打印当前时间CST格式的时间：time.now = ",time.Now().Format("2006-01-02 15:04:05"))
	filePath := time.Now().Format("20060102150405")
	//操作数据  SaveToFile：beego存储文件数据的方法    fileName：格式化时间  ext：后缀

	this.SaveToFile("FileName","./static/img/"+filePath+ext)
	return  "/static/img/"+filePath+ext  //返回文件路径
}

//处理编辑文章
func (this*ArticleController)HandleEditArticle(){
	//1获取数据
	id,err:=this.GetInt("id")//获取文章隐藏于的id
	articleName:=this.GetString("articleName")//获取文章标题
	content:=this.GetString("content")//获取文章内容
	filePath:=UploadFunc(this,"uploadname") //获取图片文件

	//2 校验数据
	if err!=nil || articleName== "" || content== "" || filePath == ""{
		beego.Error("获取数据失败",err)
		this.TplName="update.html"
		return
	}
	//3 处理数据:更新数据库数据
	// 获取orm对象
	o:=orm.NewOrm()

	//获取操作对象
	var article models.Article

	//指定操作条件
	article.Id2=id

	//更新前Read以下,判断更新的数据是否存在
	err=o.Read(&article)
	//更新
	if err!=nil{
		beego.Error("更新数据不存在")
		this.TplName="update.html"
		return
	}
	//更新
	article.Title=articleName
	article.Content=content
	article.Img=filePath
	o.Update(&article)//更新

	//4 返回数据
	this.Redirect("/article/index",302)

}

//删除文章
func (this*ArticleController)HandleDelete(){
	//获取数据
	id,err:=this.GetInt("id")

	//校验数据
	if err!=nil{
		beego.Error("删除请求数据失败",err)
		this.TplName="index.html"
		return
	}
	//处理数据
	//删除数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取删除对象
	var article models.Article
	//指定删除条件
	article.Id2=id
	//删除
	_,err=o.Delete(&article)
	if err!=nil{
		beego.Error("删除数据错误")
		this.TplName="index.html"
		return
	}
	//返回数据
	this.Redirect("/article/index",302)

}

//展示添加文章类新页面
func(this *ArticleController)ShowAddType(){
	//获取所有文章类型数据
	//获取orm对象
	o:=orm.NewOrm()
	qs:=o.QueryTable("articleType")
	var articleTypes []models.ArticleType
	qs.All(&articleTypes)
	//传递数据给前端
	this.Data["articleTypes"]=articleTypes


	//展示layout
	this.Layout="layout.html"

	//指定页面
	this.TplName="addType.html"
}

//添加文章分类
func(this *ArticleController)HandleAddType(){
	//获取数据
	typeName:=this.GetString("typeName")

	//校验数据
	if typeName==""{
		beego.Error("添加文章类型不能为空")
		this.TplName="addType"
		return

	}
	//处理数据
	//插入
	//获取orm对象
	o:=orm.NewOrm()
	//获取插入对象
	var articleType models.ArticleType
	//给插入对象赋值
	articleType.TypeName=typeName
	//插入
	o.Insert(&articleType)
	//返回数据
	this.Redirect("/article/addType",302)
}

//删除文章类型
func(this *ArticleController)DeleteType(){
	//获取数据
	id,err:=this.GetInt("id")

	//校验数据
	if err!=nil {
		beego.Error("获取类型数据失败")
		this.TplName="addType.html"
		return
	}

	//处理数据  删除数据
	//获取orm对象
	o:=orm.NewOrm()
	//获取删除对象
	var articleType models.ArticleType
	//指定条件
	articleType.Id=id
	//删除
	o.Delete(&articleType)
	//返回数据
	this.Redirect("/article/addType",302)
}
