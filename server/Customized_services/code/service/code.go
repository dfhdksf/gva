package service

import (
	"context"
	"fmt"
	smsv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/sms/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/Customized_services/code/repository"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

const codeTplId = "1877556"

//go:generate mockgen -source=./code.go -package=svcmocks -destination=mocks/code.mock.go CodeService
type CodeService interface {
	//每个业务场景的验证码是独立的，用biz区分业务
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type SMSCodeService struct {
	sms  smsv1.SmsServiceClient //sms指短信服务，验证码集成短信服务
	repo repository.CodeRepository
	//tplId string 这里可以让别人配置，但这种短信验证服务可能很久模板都不会变，所以也可以尝试定义为常量
}

func NewSMSCodeService(svc smsv1.SmsServiceClient, repo repository.CodeRepository) CodeService {
	return &SMSCodeService{
		sms:  svc,
		repo: repo,
	}
}

// Send 生成一个随机验证码，并发送
func (c *SMSCodeService) Send(ctx context.Context, biz string, phone string) error {
	//生成验证码
	code := c.generate()
	//存储到redis中
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	_, err = c.sms.Send(ctx, &smsv1.SmsSendRequest{
		TplId:   codeTplId,
		Args:    []string{code},
		Numbers: []string{phone},
	})
	//思考：如果err!=nil，那么redis有验证码用户不一定接收到验证码，要不要将redis中验证码删除
	//答：不能，因为err可能为超时的err，都不知道验证码发送出去没，发送出去的话你删掉就更麻烦了
	//可以在这里重试，配置一个重试的服务
	return err
}

// Verify 验证验证码
func (c *SMSCodeService) Verify(ctx context.Context,
	biz string,
	phone string,
	inputCode string) (bool, error) {
	ok, err := c.repo.Verify(ctx, biz, phone, inputCode)
	// 这里我们在 service 层面上对 RedisHandler 屏蔽了最为特殊的错误
	if err == repository.ErrCodeVerifyTooManyTimes {
		// 在接入了告警之后，这边要告警
		// 因为这意味着有人在搞你
		return false, nil
	}
	return ok, err
}

func (c *SMSCodeService) generate() string {
	// 用随机数生成一个0-999999的数字
	num := rand.Intn(999999)
	return fmt.Sprintf("%6d", num)
}
