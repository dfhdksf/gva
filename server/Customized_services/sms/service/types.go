package service

import "context"

// Service 发送短信的抽象
// 目前你可以理解为，这是一个为了适配不同的短信供应商的抽象
//
//go:generate mockgen -source=./types.go -package=smsmocks -destination=mocks/sms.mock.go Service
type Service interface {
	Send(ctx context.Context, tplId string,
		args []string, numbers ...string) error
	//阿里云的args为map，重构接口时可以修改为args []Args，然后在具体实现里面将[]映射为map，实现了扩展
	//SendV1(ctx context.Context, tplId string,
	//	args []Args, numbers ...string) error
}

//type Args struct {
//	Name string
//	Val string
//}
