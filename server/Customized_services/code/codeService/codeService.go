package code

import (
	"context"
	"fmt"
	"math/rand"
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 短信发送接口
type SMSSender interface {
	Send(ctx context.Context, phone string, args []string, templateId string) error
}

// 验证码服务接口
type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// 验证码服务实现
type codeServiceImpl struct {
	smsSender  SMSSender
	codeStore  CodeStore
	templateId string
}

// 验证码存储接口
type CodeStore interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

var (
	ErrCodeSendTooMany        = fmt.Errorf("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = fmt.Errorf("验证次数太多")
)

// 创建验证码服务
func NewCodeService(sender SMSSender, store CodeStore) CodeService {
	return &codeServiceImpl{
		smsSender:  sender,
		codeStore:  store,
		templateId: global.GVA_CONFIG.Code.TemplateId,
	}
}

func (c *codeServiceImpl) Send(ctx context.Context, biz string, phone string) error {
	// 生成验证码
	code := c.generate()

	// 存储验证码
	err := c.codeStore.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	// 发送验证码
	return c.smsSender.Send(ctx, phone, []string{code}, c.templateId)
}

func (c *codeServiceImpl) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.codeStore.Verify(ctx, biz, phone, inputCode)
}

func (c *codeServiceImpl) generate() string {
	num := rand.Intn(999999)
	return fmt.Sprintf("%06d", num)
}
