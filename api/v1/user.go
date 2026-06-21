package v1

import (
	"aresumei/pkg/util"
	"aresumei/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, resp)
	}
}
func UserVaildEmail(c *gin.Context) {
	var sendService service.SendEmailService
	if err := c.ShouldBind(&sendService); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := sendService.Send(c.Request.Context())
		c.JSON(http.StatusOK, resp)
	}
}
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, resp)
	}
}
func UserUpLoadResume(c *gin.Context) {
	var uploadResume service.UserService
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	claims, err := util.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if err := c.ShouldBind(&uploadResume); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := uploadResume.UpLoadResume(c.Request.Context(), file, claims.ID)
		c.JSON(http.StatusOK, resp)
	}
}
func ResumeAnswer(c *gin.Context) {
	var ResumeService service.ResumeService
	claims, err := util.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if err := c.ShouldBind(&ResumeService); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := ResumeService.Resume(c.Request.Context(), claims.ID)
		c.JSON(http.StatusOK, resp)
	}
}
func UserUpLoadCompany(c *gin.Context) {
	var uploadCompany service.TextService
	claims, err := util.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	if err := c.ShouldBind(&uploadCompany); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		resp := uploadCompany.UploadCompany(c.Request.Context(), claims.ID)
		c.JSON(http.StatusOK, resp)
	}

}
