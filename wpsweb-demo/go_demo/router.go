package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)


func InitRouter(router *gin.Engine) {
	//加载静态资源
	router.StaticFS("/js",http.Dir("../app"))
	router.LoadHTMLGlob("../app/*.html")

	//文件列表页面
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	//文件展示页面
	router.GET("/view.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "view.html", gin.H{})
	})

	//前端获取dir文件夹下的文件名,仅供参考,开发者可以按照本身需求重新定义
	router.GET("/getListFile",getListFileHandler)

	//文件下载,仅供参考,开发者可以按照本身需求重新定义
	router.GET("/v1/file", CheckToken, getFileHanlder)

	//文件上传,仅供参考,开发者可以按照本身需求重新定义
	router.POST("/v1/upload", CheckToken, fileUploadHandler)

	//传入文件名,返回有效的url和token,仅供参考,开发者可以按照本身需求重新定义
	router.GET("/getUrlAndToken",getUrlAndTokenHandler)

	//生成url接口,仅提供参考,本demo未使用
	router.GET("/v1/url", getUrlHandler)



	r := router.Group("/v1/3rd")
	{
		//获取文件元数据
		r.GET("/file/info", CheckOpenSignature,CheckToken, fileHandler)

		//获取用户信息
		r.POST("/user/info", CheckOpenSignature, CheckToken,getUserBatch)

		//上传文件新版本
		r.POST("/file/save", CheckOpenSignature,CheckToken, postFile)

		//上传在线编辑用户信息
		r.POST("/file/online", CheckToken,PostFileOnline)

		//获取特定版本的文件信息
		r.GET("/file/version/:version", CheckOpenSignature,CheckToken, fileVersion)

		//文件重命名
		r.PUT("/file/rename", CheckToken,putFileName)

		//获取所有历史版本文件信息
		r.POST("/file/history", CheckToken,GetFileHistoryVersions)

		//新建文件
		r.POST("/file/new", CheckToken, postNewFile)

	}

}