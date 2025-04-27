package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

type Service struct {
	client   *sms.Client
	appId    *string //appId和signName是基本固定的signName为签发人，例如腾讯科技
	signName *string
}

func NewService(c *sms.Client, appId string,
	signName string) *Service {
	return &Service{
		client:   c,
		appId:    ekit.ToPtr[string](appId), //腾讯云短信服务要求这两个字段为*string，所以这里用泛型工具库将string转为指针(return &string)
		signName: ekit.ToPtr[string](signName),
	}
}

//tplId：短信模板 ID，标识所使用的短信模板，决定短信的格式内容，在腾讯云短信后台配置，需要和具体模板一一对应
//args：短信模板参数，用于动态填充模板内容中的占位符 (如验证码、用户姓名等)
//numbers 变长参数，表示短信将要发送的目标手机号，支持一次发送给多个手机号
func (s *Service) Send(ctx context.Context, tplId string,
	args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.PhoneNumberSet = toStringPtrSlice(numbers)
	req.SmsSdkAppId = s.appId
	// ctx 继续往下传
	req.SetContext(ctx)
	req.TemplateParamSet = toStringPtrSlice(args)
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.SignName = s.signName
	//成功的返回结果保存在 `resp`，错误信息保存在 `err`
	resp, err := s.client.SendSms(req)
	zap.L().Debug("调用腾讯短信服务",
		zap.Any("req", req),
		zap.Any("resp", resp))
	if err != nil {
		return err
	}
	//检查短信发送状态，
	//遍历响应的 `SendStatusSet` 数组，检查每个手机号的发送状态。
	//如果任何状态出现错误（`Code` 不为 `Ok`），返回对应的错误
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code: %s, 原因：%s",
				*status.Code, *status.Message)
		}
	}
	return nil
}

func toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
