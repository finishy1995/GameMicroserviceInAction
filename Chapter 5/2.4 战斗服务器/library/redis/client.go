package redis

import (
	"ProjectX/library/log"
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

// Client Redis客户端
type Client struct {
	typ         string
	isLogDetail bool
	cmd         redis.Cmdable
	ctx         context.Context
}

// NewRedisClient 根据传入的Config新建Client
func NewRedisClient(config *ClientConfig) *Client {
	var cmd redis.Cmdable
	if config.Type == NodeMode {
		client := redis.NewClient(&redis.Options{
			Addr:         config.Host,
			Password:     config.Pass,
			DialTimeout:  config.ConnectTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolSize:     config.MaxActive,
			MinIdleConns: config.MinIdle,
			IdleTimeout:  config.IdleTimeout,
		})
		cmd = client
	} else {
		addr := strings.Split(config.Host, ",")
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        addr,
			Password:     config.Pass,
			DialTimeout:  config.ConnectTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolSize:     config.MaxActive,
			MinIdleConns: config.MinIdle,
			IdleTimeout:  config.IdleTimeout,
		})
		cmd = client
	}

	if config.Verbose {
		log.Info("create RedisClient config = %v", config)
	}
	return &Client{
		typ:         config.Type,
		cmd:         cmd,
		ctx:         context.Background(),
		isLogDetail: config.Verbose,
	}
}

func getExpireDuration(expire int64) time.Duration {
	return time.Second * time.Duration(expire)
}

func (rc *Client) log(format string, v ...interface{}) {
	if rc.isLogDetail {
		log.Info(format, v...)
	}
}

func (rc *Client) Type() string {
	return rc.typ
}

// Set sets a key-value pair.
func (rc *Client) Set(k string, v string) error {
	reply, err := rc.cmd.Set(rc.ctx, k, v, 0).Result()
	rc.log("set %s %s, reply=%s", k, v, reply)
	return err
}

// SetWithExpire sets a key-value pair with expire.
func (rc *Client) SetWithExpire(k string, v string, expire int64) error {
	if expire < 1 {
		return ErrInvalidExpireParameter
	}

	reply, err := rc.cmd.Set(rc.ctx, k, v, getExpireDuration(expire)).Result()
	rc.log("set %s %s, reply=%s", k, v, reply)
	return err
}

// SetNXWithExpire sets a key-value pair with expire.
// True will be returned if target pair not exists, otherwise false will be returned
func (rc *Client) SetNXWithExpire(k string, v string, expire int64) (bool, error) {
	if expire < 1 {
		return false, ErrInvalidExpireParameter
	}

	reply, err := rc.cmd.SetNX(rc.ctx, k, v, getExpireDuration(expire)).Result()
	rc.log("set %s %s ex %d nx, reply=%v", k, v, expire, reply)
	return reply, err
}

// SetExpire 修改key的过期时间
func (rc *Client) SetExpire(key string, expire int64) error {
	err := rc.cmd.Expire(rc.ctx, key, getExpireDuration(expire)).Err()
	rc.log("set %s expire %d", key, expire)
	return err
}

// Get 通过键获取对应值，键值对不存在时返回空字符串
func (rc *Client) Get(k string) (string, error) {
	reply, err := rc.cmd.Get(rc.ctx, k).Result()
	rc.log("get %s, reply=%s", k, reply)
	return reply, err
}

// Del 删除键值对
func (rc *Client) Del(k ...string) error {
	reply, err := rc.cmd.Del(rc.ctx, k...).Result()
	rc.log("del %s, reply=%d", k, reply)
	return err
}

// Expire 设置键值对超时时间
func (rc *Client) Expire(k string, expire int64) bool {
	reply, err := rc.cmd.Expire(rc.ctx, k, getExpireDuration(expire)).Result()
	rc.log("expire %s %d, reply=%v", k, expire, reply)
	if err != nil {
		return false
	}
	return reply
}

// Exists 判断键值对是否存在
func (rc *Client) Exists(k string) (bool, error) {
	reply, err := rc.cmd.Exists(rc.ctx, k).Result()
	if err != nil {
		return false, err
	}
	rc.log("exists %s, reply=%v", k, reply)
	result := true
	if reply <= 0 {
		result = false
	}
	return result, nil
}

// ListLength llen 获取列表长度
func (rc *Client) ListLength(listKey string) (int, error) {
	reply, err := rc.cmd.LLen(rc.ctx, listKey).Result()
	if err != nil {
		return -1, err
	}
	rc.log("llen %s, reply=%d", listKey, reply)
	return int(reply), nil
}

// ListRange lrange 获取列表指定范围内的元素
func (rc *Client) ListRange(listKey string, start int, stop int) ([]string, error) {
	reply, err := rc.cmd.LRange(rc.ctx, listKey, int64(start), int64(stop)).Result()
	rc.log("lrange %s %d %d, reply=%v", listKey, start, stop, reply)
	return reply, err
}

// ListTrim ltrim 对一个列表进行修改,只保留指定区间内的元素，不在指定区间之内的元素都将被删除
func (rc *Client) ListTrim(listKey string, start int, stop int) error {
	reply, err := rc.cmd.LTrim(rc.ctx, listKey, int64(start), int64(stop)).Result()
	rc.log("ltrim %s %d %d, reply=%s", listKey, start, stop, reply)
	return err
}

// ListRemove lrem 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素
func (rc *Client) ListRemove(listKey string, count int, value string) error {
	reply, err := rc.cmd.LRem(rc.ctx, listKey, int64(count), value).Result()
	rc.log("lrem %s %d %s, reply=%v", listKey, count, value, reply)
	return err
}

// RightPush rpush 将值添加到列表尾部
func (rc *Client) RightPush(listKey string, v ...string) error {
	length := len(v)
	interfaceV := make([]interface{}, length, length)
	for index, value := range v {
		interfaceV[index] = value
	}
	reply, err := rc.cmd.RPush(rc.ctx, listKey, interfaceV...).Result()
	rc.log("rpush %s %s, reply=%v", listKey, v, reply)
	return err
}

// HashSet 将哈希表 key 中的字段 field 的值设为 value
func (rc *Client) HashSet(hashKey string, field string, value string) error {
	reply, err := rc.cmd.HSet(rc.ctx, hashKey, field, value).Result()
	rc.log("hset %s %s %s, reply=%v", hashKey, field, value, reply)
	return err
}

// HashSetMap 将 map[string]interface 的内容存入对应key的哈希表中
func (rc *Client) HashSetMap(hashKey string, value map[string]interface{}) error {
	reply, err := rc.cmd.HSet(rc.ctx, hashKey, value).Result()
	rc.log("hset %s %v, reply=%v", hashKey, value, reply)
	return err
}

// HashGet 获取存储在哈希表中指定字段的值，未查找到则返回 err
// See ErrNil
func (rc *Client) HashGet(hashKey string, field string) (string, error) {
	reply, err := rc.cmd.HGet(rc.ctx, hashKey, field).Result()
	rc.log("hget %s %s, reply=%s", hashKey, field, reply)
	return reply, err
}

// HashGetAll 获取在哈希表中指定 key 的所有字段和值，不存在值则返回一个非nil的map，其len为0
func (rc *Client) HashGetAll(hashKey string) (map[string]string, error) {
	reply, err := rc.cmd.HGetAll(rc.ctx, hashKey).Result()
	rc.log("hgetall %s, reply=%s", hashKey, reply)
	return reply, err
}

// HashLen 获取在哈希表中指定Key的长度
func (rc *Client) HashLen(hashKey string) (int64, error) {
	reply, err := rc.cmd.HLen(rc.ctx, hashKey).Result()
	rc.log("hlen %s, reply=%d", hashKey, reply)
	return reply, err
}

// HashExist 获取在哈希表指定的key中是否存在目标字段
func (rc *Client) HashExist(hashKey string, field string) (bool, error) {
	reply, err := rc.cmd.HExists(rc.ctx, hashKey, field).Result()
	rc.log("hexist %s, reply=%v", hashKey, reply)
	return reply, err
}

// HashSetNX 只在值不存在时设置值
func (rc *Client) HashSetNX(hashKey string, field string, value string) (bool, error) {
	reply, err := rc.cmd.HSetNX(rc.ctx, hashKey, field, value).Result()
	rc.log("hsetnx %s %s %s, reply=%v", hashKey, field, value, reply)
	return reply, err
}

// HashDel 删除一个或多个哈希表字段
func (rc *Client) HashDel(hashKey string, field ...string) error {
	reply, err := rc.cmd.HDel(rc.ctx, hashKey, field...).Result()
	rc.log("hdel %s %s, reply=%d", hashKey, field, reply)
	return err
}
