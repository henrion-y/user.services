package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"

	"github.com/henrion-y/base.services/infra/xerror"
)

func responseError(c *gin.Context, err error) {
	if xerr, ok := err.(*xerror.XError); ok {
		log.Errorf("%+v", xerr)
		c.AbortWithStatusJSON(http.StatusOK, xerr)
	} else {
		responseError(c, xerror.NewXErrorByCode(xerror.ErrRuntime))
	}
}

type ResponseData struct {
	Code int32       `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

func responseData(c *gin.Context, data interface{}) {
	b, _ := json.Marshal(data)
	log.Debugf("%s %s response: %s", strings.ToUpper(c.Request.Method), c.Request.URL.Path, string(b))
	c.JSON(http.StatusOK, ResponseData{Code: 0, Data: data})
}

func responseSuccess(c *gin.Context) {
	log.Debugf("%s %s response success", strings.ToUpper(c.Request.Method), c.Request.URL.Path)
	c.JSON(http.StatusOK, ResponseData{Code: 0})
}
