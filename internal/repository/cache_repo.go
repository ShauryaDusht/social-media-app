package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"social-media-app/internal/config"
	"social-media-app/internal/models"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	SetTimeline(userID uint, posts []models.PostResponse, expiry time.Duration) error
	GetTimeline(userID uint) ([]models.PostResponse, error)
	DeleteTimeline(userID uint) error
	SetPostCache(postID uint, post models.PostResponse, expiry time.Duration) error
	GetPostCache(postID uint) (*models.PostResponse, error)
	DeletePostCache(postID uint) error
}

type cacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheRepository(cfg *config.Config) CacheRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &cacheRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *cacheRepository) SetTimeline(userID uint, posts []models.PostResponse, expiry time.Duration) error {
	key := getTimelineKey(userID)
	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, expiry).Err()
}

func (r *cacheRepository) GetTimeline(userID uint) ([]models.PostResponse, error) {
	key := getTimelineKey(userID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var posts []models.PostResponse
	err = json.Unmarshal([]byte(data), &posts)
	return posts, err
}

func (r *cacheRepository) DeleteTimeline(userID uint) error {
	key := getTimelineKey(userID)
	return r.client.Del(r.ctx, key).Err()
}

func (r *cacheRepository) SetPostCache(postID uint, post models.PostResponse, expiry time.Duration) error {
	key := getPostKey(postID)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, expiry).Err()
}

func (r *cacheRepository) GetPostCache(postID uint) (*models.PostResponse, error) {
	key := getPostKey(postID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var post models.PostResponse
	err = json.Unmarshal([]byte(data), &post)
	return &post, err
}

func (r *cacheRepository) DeletePostCache(postID uint) error {
	key := getPostKey(postID)
	return r.client.Del(r.ctx, key).Err()
}

func getTimelineKey(userID uint) string {
	return fmt.Sprintf("timeline:%d", userID)
}

func getPostKey(postID uint) string {
	return fmt.Sprintf("post:%d", postID)
}
