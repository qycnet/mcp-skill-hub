package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Service 缓存服务
type Service struct {
	client *redis.Client
}

// Config 缓存配置
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewService 创建缓存服务
func NewService(ctx context.Context, cfg Config) (*Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis 连接失败：%w", err)
	}

	return &Service{client: client}, nil
}

// Get 获取缓存
func (s *Service) Get(ctx context.Context, key string, value interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// Set 设置缓存
func (s *Service) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, expiration).Err()
}

// Delete 删除缓存
func (s *Service) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (s *Service) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// CacheKey 缓存键生成器
type CacheKey struct{}

// Skill 技能相关缓存键
func (CacheKey) Skill(id uint) string {
	return fmt.Sprintf("skill:%d", id)
}

// SkillList 技能列表缓存键
func (CacheKey) SkillList(page, pageSize int, category, sort string) string {
	return fmt.Sprintf("skills:list:%d:%d:%s:%s", page, pageSize, category, sort)
}

// User 用户相关缓存键
func (CacheKey) User(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

// Subscription 订阅缓存键
func (CacheKey) Subscription(userID uint) string {
	return fmt.Sprintf("subscription:%d", userID)
}

// Plan 价格计划缓存键
func (CacheKey) Plan() string {
	return "plans:all"
}

// DefaultExpiration 默认过期时间
const (
	ShortCache  = 5 * time.Minute  // 短期缓存（频繁变化数据）
	MediumCache = 30 * time.Minute // 中期缓存
	LongCache   = 24 * time.Hour   // 长期缓存（静态数据）
)

// InvalidateSkill 使技能缓存失效
func (s *Service) InvalidateSkill(ctx context.Context, id uint) error {
	key := CacheKey{}.Skill(id)
	return s.Delete(ctx, key)
}

// InvalidateUser 使用户缓存失效
func (s *Service) InvalidateUser(ctx context.Context, id uint) error {
	key := CacheKey{}.User(id)
	return s.Delete(ctx, key)
}

// InvalidatePlans 使价格计划缓存失效
func (s *Service) InvalidatePlans(ctx context.Context) error {
	key := CacheKey{}.Plan()
	return s.Delete(ctx, key)
}

// GetWithLoader 获取缓存，如果不存在则使用 loader 函数加载
func (s *Service) GetWithLoader(ctx context.Context, key string, expiration time.Duration, loader func() (interface{}, error), value interface{}) error {
	// 尝试从缓存获取
	exists, err := s.Exists(ctx, key)
	if err == nil && exists {
		return s.Get(ctx, key, value)
	}

	// 缓存未命中，加载数据
	data, err := loader()
	if err != nil {
		return err
	}

	// 存入缓存
	if err := s.Set(ctx, key, data, expiration); err != nil {
		return err
	}

	// 复制数据到输出
	dataJSON, _ := json.Marshal(data)
	return json.Unmarshal(dataJSON, value)
}
