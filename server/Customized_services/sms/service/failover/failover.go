package failover

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/sms/service"
	"sync/atomic"
)

type FailoverSMSService struct {
	// 一大堆可供你选择的 SMS Service 实现
	svcs []service.Service

	idx uint64
}

func NewFailoverSMSService(svcs []service.Service) *FailoverSMSService {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

//每次从第一个服务商开始尝试
//问题：容易造成首节点压力过大
func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		// 这边要打印日志
	}
	return errors.New("发送失败，所有服务商都尝试过了")
}

//设计模式为： 故障转移（Failover）模式 和 负载均衡策略
//通过原子操作实现轮询起点，避免单点压力，处理了上下文取消和超时的情况
func (f *FailoverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 二话不说先把下标往后推一位
	idx := atomic.AddUint64(&f.idx, 1) //原子操作将起始服务商下标后移一位
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length] //环形选择服务商
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			//	调用者设置的超时时间到了
			// 调用者主动取消了
			return err
		}
		// 其它情况会走到这里，这边要打印日志
	}
	return errors.New("发送失败，所有服务商都尝试过了")
}
