package code

import (
	"context"
	"fmt"
)

// 默认的短信发送器实现
type defaultSMSSender struct{}

func NewDefaultSMSSender() SMSSender {
	return &defaultSMSSender{}
}

func (s *defaultSMSSender) Send(ctx context.Context, phone string, args []string, templateId string) error {
	// 这里可以集成你自己的短信发送逻辑
	fmt.Printf("发送短信到 %s, 模板ID: %s, 参数: %v\n", phone, templateId, args)
	return nil
}
