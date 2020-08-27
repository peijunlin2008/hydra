package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro-plat/hydra/components/pkgs/cache"
	"github.com/micro-plat/hydra/components/pkgs/redis"
)

// Client redis配置文件
type Client struct {
	servers []string
	client  *redis.Client
}

// New 根据配置文件创建一个redis连接
func New(addrs []string, raw string) (m *Client, err error) {
	m = &Client{servers: addrs}
	m.client, err = redis.New(redis.WithRaw(raw))
	if err != nil {
		return
	}
	m.servers = m.client.GetAddrs()
	return
}

//GetServers 获取服务器列表
func (c *Client) GetServers() []string {
	return c.servers
}

// Get 根据key获取redis中的数据
func (c *Client) Get(key string) (string, error) {
	data, err := c.client.Get(key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return data, nil
		}
		return "", err
	}
	return data, nil
}

//Decrement 增加变量的值
func (c *Client) Decrement(key string, delta int64) (n int64, err error) {
	return c.client.DecrBy(key, delta).Result()
}

//Increment 减少变量的值
func (c *Client) Increment(key string, delta int64) (n int64, err error) {
	return c.client.IncrBy(key, delta).Result()
}

//Gets 获取多条数据
func (c *Client) Gets(key ...string) (r []string, err error) {
	data, err := c.client.MGet(key...).Result()
	if err != nil {
		return nil, err
	}
	r = make([]string, 0, len(data))
	for _, v := range data {
		if v == nil || v.(string) == "" {
			continue
		}
		r = append(r, v.(string))
	}
	return
}

// Add 添加数据到redis中,如果redis存在，则报错
func (c *Client) Add(key string, value string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	i, err := c.client.Exists(key).Result()
	if err != nil {
		return err
	}
	if i == 1 {
		err = fmt.Errorf("key:%s已存在", key)
		return err
	}
	_, err = c.client.Set(key, value, expires).Result()
	return err
}

// Set 更新数据到redis中，没有则添加
func (c *Client) Set(key string, value string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	_, err := c.client.Set(key, value, expires).Result()
	return err
}

//Delete 删除指定的KEY,支持*模糊匹配
func (c *Client) Delete(key string) error {
	if !strings.Contains(key, "*") {
		_, err := c.client.Del(key).Result()
		if err != nil {
			return fmt.Errorf("%v(%s)", err, key)
		}
		return nil
	}
	_, err := c.client.Eval(`
    local keys=redis.call('KEYS',KEYS[1])
    if (#keys==0) then
        return 0
    end
		return redis.call('DEL',unpack(keys))`, []string{key}).Result()
	return err
}

//Exists 查询指定的KEY是否存在
func (c *Client) Exists(key string) bool {
	r, err := c.client.Exists(key).Result()
	return err == nil && r == 1
}

//Delay 延长数据在redis中的时间
func (c *Client) Delay(key string, expiresAt int) error {
	expires := time.Duration(expiresAt) * time.Second
	if expiresAt == 0 {
		expires = 0
	}
	_, err := c.client.Expire(key, expires).Result()
	return err
}

//Close 关闭服务器连接
func (c *Client) Close() error {
	return c.client.Close()
}

type redisResolver struct {
}

func (s *redisResolver) Resolve(address []string, conf string) (cache.ICache, error) {
	return New(address, conf)
}
func init() {
	cache.Register("redis", &redisResolver{})
}
