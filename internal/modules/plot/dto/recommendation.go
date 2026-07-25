package dto

type BedRecommendationsResponse struct {
	Recommendations []CropRecommendation `json:"recommendations"`
	SearchResults   []CropSearchResult   `json:"search_results,omitempty"`
}

type CropRecommendation struct {
	CropID     int32          `json:"crop_id"`
	Name       string         `json:"name"`
	FamilyName string         `json:"family_name"`
	Score      int32          `json:"-"`
	Reasons    []ReasonDetail `json:"reasons"`
}

type ReasonDetail struct {
	Explanation string `json:"explanation"`
	IsPositive  bool   `json:"ispositive"`
	Score       int32  `json:"-"`
}

type CropSearchResult struct {
	CropID     int32  `json:"crop_id"`
	Name       string `json:"name"`
	FamilyName string `json:"family_name"`
}
