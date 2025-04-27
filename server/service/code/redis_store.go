package code

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

// RedisCodeStore Redis 验证码存储
type RedisCodeStore struct {
	redis          redis.Cmdable
	expiration     time.Duration
	resendInterval time.Duration
}

// NewRedisCodeStore 创建Redis验证码存储
func NewRedisCodeStore(redis redis.Cmdable, expiration, resendInterval int) CodeStore {
	return &RedisCodeStore{
		redis:          redis,
		expiration:     time.Duration(expiration) * time.Second,
		resendInterval: time.Duration(resendInterval) * time.Second,
	}
}

func (r *RedisCodeStore) Store(ctx context.Context, biz string, phone string, code string) error {
	key := r.key(biz, phone)
	res, err := r.redis.Eval(ctx, luaSetCode, []string{key}, code, int(r.expiration.Seconds()), int(r.resendInterval.Seconds())).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return fmt.Errorf("未知错误码: %d", res)
	}
}

func (r *RedisCodeStore) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	key := r.key(biz, phone)
	res, err := r.redis.Eval(ctx, luaVerifyCode, []string{key}, inputCode).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooManyTimes
	default:
		return false, nil
	}
}

func (r *RedisCodeStore) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
