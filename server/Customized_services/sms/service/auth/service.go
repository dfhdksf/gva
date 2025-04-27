package auth

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/sms/service"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc service.Service
	key []byte
}

func (s *SMSService) Send(ctx context.Context,
// 改变了语义
	tplToken string, args []string, numbers ...string) error {
	var c Claims
	//ParseWithClaims负责 解析、验证 JWT 令牌并提取声明（Claims）数据。具体实现以下关键功能：
	//将字符串形式的 tplToken 解析为 JWT 的三个组成部分：
	//使用提供的密钥（s.key）验证签名有效性，确保令牌未被篡改
	//将解码后的声明数据填充到 Claims 结构体中
	_, err := jwt.ParseWithClaims(tplToken, &c, func(token *jwt.Token) (interface{}, error) { //生成的token暴露给客户端后客户端携带上token发送请求然后服务端在这里解析，通过就发送
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, c.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}

/*
颁发tplToken:
// 授权服务示例代码
func GenerateTplToken(userID string, tplID string) (string, error) {
    claims := auth.Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
            Subject:   userID,
        },
        Tpl: tplID,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(sharedSecretKey) // 与SMSService使用相同密钥
}

*/
