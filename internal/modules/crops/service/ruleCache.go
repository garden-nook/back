package service

import (
	"context"
	"garden-nook/internal/modules/crops/model"
	"garden-nook/internal/modules/crops/repository"
	"sync"
)

type RuleCache struct {
	mu       sync.RWMutex
	byCrop   map[int32][]model.RuleInfo
	byFamily map[int32][]model.RuleInfo
	repo     *repository.CropRuleRepo
}

func NewRuleCache(repo *repository.CropRuleRepo) (*RuleCache, error) {
	c := &RuleCache{repo: repo}
	if err := c.refresh(context.Background()); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *RuleCache) GetRulesBySubjectCropID(ctx context.Context, cropID int32) ([]model.RuleInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.byCrop[cropID], nil
}

func (c *RuleCache) GetRulesBySubjectFamilyID(ctx context.Context, familyID int32) ([]model.RuleInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.byFamily[familyID], nil
}

func (c *RuleCache) Refresh(ctx context.Context) error {
	return c.refresh(ctx)
}

func (c *RuleCache) refresh(ctx context.Context) error {
	rules, err := c.repo.ListRulesFull(ctx)
	if err != nil {
		return err
	}

	byCrop := make(map[int32][]model.RuleInfo)
	byFamily := make(map[int32][]model.RuleInfo)

	for _, r := range rules {
		if r.SubjectCropID != nil {
			byCrop[*r.SubjectCropID] = append(byCrop[*r.SubjectCropID], r)
		}
		if r.SubjectFamilyID != nil {
			byFamily[*r.SubjectFamilyID] = append(byFamily[*r.SubjectFamilyID], r)
		}
	}

	c.mu.Lock()
	c.byCrop = byCrop
	c.byFamily = byFamily
	c.mu.Unlock()
	return nil
}
