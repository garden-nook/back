package model

import (
	"garden-nook/internal/modules/plot/enum"
	"time"
)

type CropInfo struct {
	ID         int32
	Name       string
	FamilyID   int32
	FamilyName string
	SunNeeds   enum.ShadeLevel
	SoilTypeID int32
}

type CropFilter struct {
	SoilTypeID *int32
	SunNeeds   []int32
	Search     string
	Limit      int
}

type PredecessorInfo struct {
	CropID          int32
	FamilyID        int32
	LastHarvestDate time.Time
}
