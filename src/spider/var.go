package spider

import (
	"github.com/garyburd/redigo/redis"
	"sync/atomic"
)

type TGlobalVars struct {
	TQueue    *TQueue
	RedisPool *redis.Pool
	index     int64
}

func (obj *TGlobalVars) Init() {
	if len(GTSpiderConfig.RedisServerAddress) > 0 {
		obj.RedisPool = NewRedisPool(GTSpiderConfig.RedisServerAddress)
	}

	go func() {
		obj.TQueue.PullToAsync(func(item TQueueItem) {
			content := item.Content()
			if obj.RedisPool != nil {
				obj.RedisPool.Get().Do("rpush", "proxy", content)
			}
			Logger.Printf("sync to redis %s", content)
		})
	}()
}

func (obj *TGlobalVars) Index() int64 {
	return atomic.AddInt64(&obj.index, 1)
}

func (obj *TGlobalVars) Clear() {
	obj.TQueue.Stop()

	if obj.RedisPool != nil {
		obj.RedisPool.Close()
	}

	Logger.Print("global vars clear exit")
}

func NewTGlobalVars() *TGlobalVars {
	return &TGlobalVars{
		TQueue : NewTQueue(1024),
		RedisPool:nil,
		index:0,
	}
}
