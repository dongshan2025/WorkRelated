package wscache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type WsRedis struct {
	Client *redis.Client
}

var redisClient *WsRedis

func GetInstance() *WsRedis {
	return redisClient
}

var once sync.Once

func NewWsRedis(addr string, password string, db int) *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	once.Do(func() {
		redisClient = &WsRedis{
			Client: redis,
		}
	})

	return redis
}

// ========== String操作 ==========
// 设置key的值
func (wr *WsRedis) Set(key string, value string) (bool, error) {
	result, err := wr.Client.Set(context.Background(), key, value, 0).Result()
	if err != nil {
		return false, err
	}

	return result == "OK", nil
}

// 设置key的值，并指定过期时间
func (wr *WsRedis) SetEx(key string, value string, ex time.Duration) (bool, error) {
	result, err := wr.Client.Set(context.Background(), key, value, ex).Result()
	if err != nil {
		return false, err
	}

	return result == "OK", nil
}

// 获取key的值
func (wr *WsRedis) Get(key string) (string, error) {
	result, err := wr.Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

// 获取旧值，同时设置新值
func (wr *WsRedis) GetSet(key string, value string) (string, error) {
	oldValue, err := wr.Client.GetSet(context.Background(), key, value).Result()
	if err != nil {
		return "", err
	}

	return oldValue, nil
}

// 指定key的值加1，并返回增加之后的值
func (wr *WsRedis) Incr(key string) (int64, error) {
	value, err := wr.Client.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return value, nil
}

// 指定key的值减1，并返回减少之后的值
func (wr *WsRedis) Decr(key string) (int64, error) {
	value, err := wr.Client.Decr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return value, nil
}

// 指定key的值减去一个指定的值，并返回减少之后的值
func (wr *WsRedis) DecrBy(key string, incr int64) (int64, error) {
	value, err := wr.Client.DecrBy(context.Background(), key, incr).Result()
	if err != nil {
		return 0, err
	}

	return value, nil
}

// 指定key的值增加一个指定的浮点数，并返回增加之后的值
func (wr *WsRedis) IncrFloat(key string, incr float64) (float64, error) {
	value, err := wr.Client.IncrByFloat(context.Background(), key, incr).Result()
	if err != nil {
		return 0, err
	}

	return value, nil
}

// 删除指定的key
func (wr *WsRedis) Del(key string) error {
	_, err := wr.Client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return nil
}

// 设置指定key的过期时间
func (wr *WsRedis) Expire(key string, ex time.Duration) (bool, error) {
	result, err := wr.Client.Expire(context.Background(), key, ex).Result()
	if err != nil {
		return false, nil
	}

	return result, nil
}

// ========== List操作 ==========
// 从列表左边插入数据，并返回列表长度
func (wr *WsRedis) LPush(key string, data ...interface{}) (int64, error) {
	result, err := wr.Client.LPush(context.Background(), key, data).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// 从列表右边插入数据，并返回列表长度
func (wr *WsRedis) RPush(key string, data ...interface{}) (int64, error) {
	result, err := wr.Client.RPush(context.Background(), key, data).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// 从列表左边删除第一个数据，并返回删除的数据
func (wr *WsRedis) LPop(key string) (string, error) {
	value, err := wr.Client.LPop(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// 从列表右边删除第一个数据，并返回删除的数据
func (wr *WsRedis) RPop(key string) (string, error) {
	value, err := wr.Client.RPop(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// 根据索引坐标，查询列表中的数据
func (wr *WsRedis) LIndex(key string, index int64) (string, error) {
	value, err := wr.Client.LIndex(context.Background(), key, index).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// 返回列表长度
func (wr *WsRedis) LLen(key string) (int64, error) {
	value, err := wr.Client.LLen(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	return value, nil
}

// 返回列表的一个范围内的数据，也可以返回全部数据
func (wr *WsRedis) LRange(key string, start, stop int64) ([]string, error) {
	values, err := wr.Client.LRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// 从列表左边开始，删除元素data，如果出现重复元素，仅删除count次
func (wr *WsRedis) LRem(key string, count int64, data interface{}) (bool, error) {
	_, err := wr.Client.LRem(context.Background(), key, count, data).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 在列表中pivot元素的后面插入data
func (wr *WsRedis) LInsert(key string, pivot int64, data interface{}) (bool, error) {
	err := wr.Client.LInsert(context.Background(), key, "after", pivot, data).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// ========== Set操作 ==========
// 添加元素到集合中，并返回集合元素个数
func (wr *WsRedis) SAdd(key string, data ...interface{}) (int64, error) {
	value, err := wr.Client.SAdd(context.Background(), key, data).Result()
	if err != nil {
		return 0, err
	}

	return value, nil
}

// 获取集合元素个数
func (wr *WsRedis) SCard(key string) (int64, error) {
	count, err := wr.Client.SCard(context.Background(), "key").Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 判断元素是否在集合中
func (wr *WsRedis) SIsMember(key string, data interface{}) (bool, error) {
	isExist, err := wr.Client.SIsMember(context.Background(), key, data).Result()
	if err != nil {
		return false, err
	}

	return isExist, nil
}

// 获取集合所有元素
func (wr *WsRedis) SMembers(key string) ([]string, error) {
	values, err := wr.Client.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return values, nil
}

// 删除key集合中的data元素
// 返回值：从集合中删除的成员数量，不包括不存在的成员
func (wr *WsRedis) SRem(key string, data ...interface{}) (int64, error) {
	count, err := wr.Client.SRem(context.Background(), key, data).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 随机返回集合中的count个元素，并且删除这些元素
func (wr *WsRedis) SPopN(key string, count int64) ([]string, error) {
	values, err := wr.Client.SPopN(context.Background(), key, count).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// ========== Hash操作 ==========
// 根据key和field字段设置field字段的值
// 返回 1 表示字段是哈希中的一个新的字段并设置字段值
// 返回 0 表示哈希中已存在此字段并更新字段值
func (wr *WsRedis) HSet(key, field, value string) (int64, error) {
	status, err := wr.Client.HSet(context.Background(), key, field, value).Result()
	if err != nil {
		return 0, err
	}

	return status, nil
}

// 根据key和field字段查询field字段的值
func (wr *WsRedis) HGet(key, field string) (string, error) {
	value, err := wr.Client.HGet(context.Background(), key, field).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

// 根据key和多个字段名，批量查询多个hash字段值
func (wr *WsRedis) HMGet(key string, fields ...string) ([]interface{}, error) {
	values, err := wr.Client.HMGet(context.Background(), key, fields...).Result()
	if err != nil {
		return nil, err
	}

	return values, nil
}

// 根据key查询所有字段和值
func (wr *WsRedis) HGetAll(key string) (map[string]string, error) {
	data, err := wr.Client.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 根据key返回所有字段名
func (wr *WsRedis) HKeys(key string) ([]string, error) {
	fields, err := wr.Client.HKeys(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// 根据key查询hash的字段数量
func (wr *WsRedis) HLen(key string) (int64, error) {
	count, err := wr.Client.HLen(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 根据key和多个字段名和字段值，批量设置hash字段值
func (wr *WsRedis) HMSet(key string, data map[string]interface{}) (bool, error) {
	result, err := wr.Client.HMSet(context.Background(), key, data).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// 如果field字段不存在，则设置hash字段值
func (wr *WsRedis) HSetNX(key, field string, value interface{}) (bool, error) {
	result, err := wr.Client.HSetNX(context.Background(), key, field, value).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// 根据key和字段名，删除hash字段，支持批量删除
func (wr *WsRedis) HDel(key string, fields ...string) (int64, error) {
	count, err := wr.Client.HDel(context.Background(), key, fields...).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 检测hash字段名是否存在
func (wr *WsRedis) HExists(key, field string) (bool, error) {
	result, err := wr.Client.HExists(context.Background(), key, field).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// ========== Zset操作 ==========
