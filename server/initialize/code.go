package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service/code"
)

func InitCodeService() {
	// 创建短信发送器
	smsSender := code.NewDefaultSMSSender()

	// 创建验证码存储
	codeStore := code.NewRedisCodeStore(
		global.GVA_REDIS,
		global.GVA_CONFIG.Code.Expiration,
		global.GVA_CONFIG.Code.ResendInterval,
	)

	// 创建验证码服务
	global.GVA_CODE = code.NewCodeService(
		smsSender,
		codeStore,
		global.GVA_CONFIG.Code.TemplateId,
	)
}
