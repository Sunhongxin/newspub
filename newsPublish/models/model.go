package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
//用户表
type User struct {
	Id int         //用户ID主键
	Name string `orm:"unique"` //用户名全局唯一
	Pwd string
	//用户和文章是多对多
	Articles []*Article `orm:"rel(m2m)"`   //多了一张表 user_article
}
//beego规定 ：当没有设置主键时，以字段名为Id，类型为int的字段 当 默认主键

//文章表
type Article struct {
	Id2 int  `orm:"pk;auto"`   //设置主键属性  主键自增
	Title string `orm:"size(40)"`//标题  标题长度
	Content string `orm:"type(text)"`//内容  内容大小
	ReadCount int 	`orm:"default(0)"`//阅读量 默认值为0 （没有阅读量）
	Time time.Time	`orm:"type(datetime);auto_now_add"`//添加时间 创建时间——更新时间
	Img string	`orm:"null"`//图片 存储路径 ：字符串   beego默认非空  空：需要手动设置
	//文章与文章类型是多对一  (类型设置外建)
	ArticleType *ArticleType `orm:"rel(fk)";on_delete(set_null);null`    //多了一个外建
	//用户
	Users  []*User `orm:"reverse(many)"`
}
type ArticleType struct {
	Id int
	TypeName string
	//文章类型和文章是一对多   文章是多
	Articles []*Article `orm:"reverse(many)"`
}


func init(){

	// orm框架的作用 1、通过orm可以创建关系型数据库表 2、通过orm操作关系型数据库表

	//注册数据库
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsPublish?charset=utf8")

	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//跑起来  false:强制更新
	orm.RunSyncdb("default",true,true)
}
