package routes

import (
	api "aresumei/api/v1"
	"aresumei/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	v1 := r.Group("/api/v1")
	{
		//用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		v1.POST("upload/resume", api.UserUpLoadResume)
		v1.POST("upload/company", api.UserUpLoadCompany)
		v1.POST("user/vaild-email", api.UserVaildEmail)
		v1.POST("user/resume-answer", api.ResumeAnswer)
		v1.POST("user/resume-optimize", api.ResumeOptimize)
	}
	return r
}
