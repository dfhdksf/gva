package code

import (
	"context"
	"fmt"
	"math/rand"

	// 使用接口包中的定义
	"github.com/flipped-aurora/gin-vue-admin/server/interfaces/code"
)

// 确保 CodeServiceImpl 实现了 interfaces/code.CodeService 接口
var _ code.CodeService = (*CodeServiceImpl)(nil)

// 短信发送接口
type SMSSender interface {
	Send(ctx context.Context, phone string, args []string, templateId string) error
}

// 验证码存储接口
type CodeStore interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

// CodeServiceImpl 验证码服务实现
type CodeServiceImpl struct {
	smsSender  SMSSender
	codeStore  CodeStore
	templateId string
}

var (
	ErrCodeSendTooMany        = fmt.Errorf("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = fmt.Errorf("验证次数太多")
)

// NewCodeService 创建验证码服务
func NewCodeService(sender SMSSender, store CodeStore, templateId string) code.CodeService {
	return &CodeServiceImpl{
		smsSender:  sender,
		codeStore:  store,
		templateId: templateId,
	}
}

func (c *CodeServiceImpl) Send(ctx context.Context, biz string, phone string) error {
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

func (c *CodeServiceImpl) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return c.codeStore.Verify(ctx, biz, phone, inputCode)
}

func (c *CodeServiceImpl) generate() string {
	num := rand.Intn(999999)
	return fmt.Sprintf("%06d", num)
}
