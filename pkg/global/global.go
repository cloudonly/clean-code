package global

import (
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/luomu/clean-code/pkg/config"
)

//var (
//	DB      *gorm.DB
//	DB_LIST map[string]*gorm.DB
//	REDIS   *redis.Client
//	CONFIG  config.Server
//	VIPER   *viper.Viper
//	LOG     *zap.Logger
//	//GLOBAL_Concurrency_Control = &singleflight.Group{}
//
//	lock sync.RWMutex
//)

var Global = &Manager{}

type Manager struct {
	lock   sync.RWMutex
	Log    *zap.Logger
	Viper  *viper.Viper
	Config config.Server
	Redis  *redis.Client
	Dbs    map[string]*gorm.DB
}

func (m *Manager) GetGlobalDBByName(dbName string) *gorm.DB {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.Dbs[dbName]
}

func (m *Manager) MustGetGlobalDBByName(dbName string) *gorm.DB {
	m.lock.RLock()
	defer m.lock.RUnlock()
	db, ok := m.Dbs[dbName]
	if !ok || db == nil {
		panic("Database init failed.")
	}
	return db
}
