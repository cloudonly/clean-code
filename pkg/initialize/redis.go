package initialize

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/luomu/clean-code/pkg/global"
)

func Redis() {
	redisConf := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password, // no password set
		DB:       redisConf.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.LOG.Error("Connect redis ping failed, err:", zap.Error(err))
	} else {
		global.LOG.Info("Connect redis ping response:", zap.String("pong", pong))
		global.REDIS = client
	}
}
