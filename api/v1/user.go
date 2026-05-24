package v1

import (
	"aresumei/pkg/util"
	"aresumei/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err != nil {

	}
}
func UserUpLoadResume(c *gin.Context) {
	resp := util.UploadFile(c)
	c.JSON(http.StatusOK, resp)
}
