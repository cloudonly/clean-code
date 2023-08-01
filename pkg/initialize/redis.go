package initialize

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/luomu/clean-code/pkg/global"
)

func Redis() {
	redisConf := global.Global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password, // no password set
		DB:       redisConf.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Global.Log.Error("Connect redis ping failed, err:", zap.Error(err))
	} else {
		global.Global.Log.Info("Connect redis ping response:", zap.String("pong", pong))
		global.Global.Redis = client
	}
}
