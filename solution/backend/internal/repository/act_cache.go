package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"neJok/solution/internal/model"
	"strconv"
	"time"
)

type ActCacheRepo struct {
	redis *redis.Client
	ctx   context.Context
}

func NewActCacheRepo(redis *redis.Client) *ActCacheRepo {
	ctx := context.Background()
	return &ActCacheRepo{redis, ctx}
}

func (r *ActCacheRepo) GetInt(key string) (int, error) {
	val, err := r.redis.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	day, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return day, nil
}

func (r *ActCacheRepo) SetInt(key string, value int) error {
	err := r.redis.Set(r.ctx, key, strconv.Itoa(value), 0).Err()
	return err
}

func (r *ActCacheRepo) GetList(key string) ([]model.CampaignForUser, error) {
	vals, err := r.redis.LRange(r.ctx, key, 0, -1).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var campaigns []model.CampaignForUser
	for _, v := range vals {
		var campaign model.CampaignForUser
		err := json.Unmarshal([]byte(v), &campaign)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (r *ActCacheRepo) SetList(key string, campaigns []model.CampaignForUser) error {
	var vals []string
	for _, campaign := range campaigns {
		data, err := json.Marshal(campaign)
		if err != nil {
			return err
		}
		vals = append(vals, string(data))
	}

	err := r.redis.Del(r.ctx, key).Err()
	if err != nil {
		return err
	}

	if len(vals) > 0 {
		err = r.redis.RPush(r.ctx, key, vals).Err()
	}
	return err
}

func (r *ActCacheRepo) DeleteKeysByPrefix(prefix string) error {
	var cursor uint64
	for {
		keys, newCursor, err := r.redis.Scan(r.ctx, cursor, prefix+"*", 0).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.redis.Del(r.ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if newCursor == 0 {
			break
		}
		cursor = newCursor
	}
	return nil
}

func (r *ActCacheRepo) SetStr(key string, value string, expire *time.Duration) error {
	var err error
	if expire != nil {
		err = r.redis.Set(r.ctx, key, value, *expire).Err()
	} else {
		err = r.redis.Set(r.ctx, key, value, 0).Err()
	}
	return err
}

func (r *ActCacheRepo) GetStr(key string) (string, error) {
	val, err := r.redis.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}
