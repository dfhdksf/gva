package failover

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/sms/service"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	//lock sync.Mutex
	svcs []service.Service
	idx  int32

	// 连续超时次数
	cnt int32

	// 连续超时次数阈值
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []service.Service, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	cnt := atomic.LoadInt32(&t.cnt)
	idx := atomic.LoadInt32(&t.idx)
	if cnt >= t.threshold {
		// 触发切换，计算新的下标
		newIdx := (idx + 1) % int32(len(t.svcs))
		// CAS 操作失败，说明有人切换了，所以你这里不需要检测返回值
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) { //将t.idx的值与idx的值进行比较，如果相等，则将newIdx赋值给t.idx，否则不进行任何操作
			//原子操作函数，用于在并发环境下安全地更新一个 int32 类型的变量。它的作用是实现 无锁的原子性比较并更新。
			//这个操作是原子的（atomic），意味着它在多线程/协程环境下是安全的，不会因为线程切换或并发竞争导致数据不一致。
			// 进入这个分支说明结果为true，当前的idx与要切换的idx相同，说明还没切换，将新的idx赋值给idx
			atomic.StoreInt32(&t.cnt, 0) //cnt赋值为0
		}
		idx = newIdx
	}
	svc := t.svcs[idx]
	// 当前使用的 svc
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case nil:
		// 没有任何错误，重置计数器
		atomic.StoreInt32(&t.cnt, 0)
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		// 如果是别的异常的话，我们保持不动
	}
	return err
}

/*
sync/atomic包能够实现原子操作
实现无锁并发安全更新：如计数器、状态变量、索引、状态机等。
CAS（Compare-And-Swap）机制的基础：是乐观并发控制的核心思想。
避免加锁带来的性能开销：尤其在读多写少的场景中非常有用。
*/
