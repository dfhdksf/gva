package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type CodeApi struct{}

func (ca *CodeApi) SendCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Biz   string `json:"biz" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数绑定失败", c)
		return
	}

	err := global.GVA_CODE.Send(c, req.Biz, req.Phone)
	if err != nil {
		response.FailWithMessage("发送验证码失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("发送验证码成功", c)
}

func (ca *CodeApi) VerifyCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Biz   string `json:"biz" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数绑定失败", c)
		return
	}

	ok, err := global.GVA_CODE.Verify(c, req.Biz, req.Phone, req.Code)
	if err != nil {
		response.FailWithMessage("验证码验证失败: "+err.Error(), c)
		return
	}

	if !ok {
		response.FailWithMessage("验证码错误", c)
		return
	}

	response.OkWithMessage("验证码验证成功", c)
}
