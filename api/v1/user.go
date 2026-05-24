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
func UserUpLoadResume(c *gin.Context) {
	resp := util.UploadFile(c)
	c.JSON(http.StatusOK, resp)
}
