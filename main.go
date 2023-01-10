package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

var (
	DB *gorm.DB
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func InitDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/acbb?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func main() {
	//连接数据库
	r := gin.Default()
	err := InitDB()
	if err != nil {
		fmt.Println("InitDB error:", err)
		return
	}

	defer DB.Close()
	DB.AutoMigrate(&Todo{})
	r.LoadHTMLGlob("./*.html")
	r.Static("/static", "static")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	v1Group := r.Group("v1")
	{
		//增
		v1Group.POST("/todo", func(c *gin.Context) {
			//存入前端返回数据
			//放入数据库中
			//返回成功状态
			var todo Todo
			c.BindJSON(&todo)
			if err := DB.Create(&todo).Error; err != nil {
				c.JSON(200, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(200, gin.H{
					"err": "无效id",
				})
				return
			}
			var todo Todo
			if err := DB.Where("id=?", id).First(&todo).Error; err != nil {
				c.JSON(200, gin.H{
					"err": err.Error(),
				})
				return
			}
			c.ShouldBindJSON(&todo)
			if err := DB.Save(&todo).Error; err != nil {
				c.JSON(200, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//查多个
		v1Group.GET("/todo", func(c *gin.Context) {
			var TodoList []Todo
			if err := DB.Find(&TodoList).Error; err != nil {
				c.JSON(200, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, TodoList)
			}
		})
		//查一个
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//删除某一个
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, _ := c.Params.Get("id")
			var todo Todo
			DB.Where("id=?", id).First(&todo)
			DB.Delete(&todo)
		})
	}
	r.Run()
}
