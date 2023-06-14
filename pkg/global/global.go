package global

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/luomu/clean-code/pkg/config"
)

var (
	DB      *gorm.DB
	DB_LIST map[string]*gorm.DB
	REDIS   *redis.Client
	CONFIG  config.Server
	VIPER   *viper.Viper
	LOG     *zap.Logger
	//GLOBAL_Concurrency_Control = &singleflight.Group{}

	lock sync.RWMutex
)

func GetGlobalDBByName(dbName string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return DB_LIST[dbName]
}

func MustGetGlobalDBByName(dbName string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := DB_LIST[dbName]
	if !ok || db == nil {
		panic("Database init failed.")
	}
	return db
}
