package crops

// SoilType - доменная модель типа почвы
type SoilType struct {
	ID          int32   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
}

// CropFamily — доменная модель семейства культур.
type CropFamily struct {
	ID          int32   `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
}

// Crop — доменная модель культуры.
type Crop struct {
	ID                int32    `json:"id" db:"id"`
	Name              string   `json:"name" db:"name"`
	Description       *string  `json:"description,omitempty" db:"description"`
	FamilyID          int32    `json:"family_id" db:"family_id"`
	FamilyName        *string  `json:"family_name,omitempty" db:"family_name"`
	VegetationDaysAvg int32    `json:"vegetation_days_avg" db:"vegetation_days_avg"`
	SunNeeds          SunNeeds `json:"sun_needs" db:"sun_needs"`
	SoilTypeID        int32    `json:"soil_type_id" db:"soil_type_id"`
	SoilName          *string  `json:"soil_name,omitempty" db:"soil_name"`
	IsDeleted         bool     `json:"-" db:"is_deleted"`
}

// CropExtended — доменная модель культуры (содержит информацию о совместимости).
type CropExtended struct {
	*Crop          `json:"crop"`
	*CropRelations `json:"crop_relations"`
}

// CropRule — правило совместимости (севооборот/аллелопатия).
type CropRule struct {
	ID              int32       `json:"id" db:"rule_id"`
	SubjectCropID   *int32      `json:"subject_crop_id,omitempty" db:"subject_crop_id"`
	SubjectFamilyID *int32      `json:"subject_family_id,omitempty" db:"subject_family_id"`
	ContextType     ContextType `json:"context_type" db:"context_type"`
	ContextCropID   *int32      `json:"context_crop_id,omitempty" db:"context_crop_id"`
	ContextFamilyID *int32      `json:"context_family_id,omitempty" db:"context_family_id"`
	ReturnAfterDays int32       `json:"return_after_days" db:"return_after_days"`
	ScoreModifier   int32       `json:"score_modifier" db:"score_modifier"`
	Explanation     string      `json:"explanation" db:"explanation"`
	Priority        int32       `json:"priority" db:"priority"`
}

// CropRelation – одна связь культуры с соседом/предшественником/последователем.
type CropRelation struct {
	CropID   int32  `json:"crop_id"`
	CropName string `json:"crop_name"`
	Score    int32  `json:"-"`
	//Explanation string `json:"explanation"`
}

// FamilyRelation – связь с семейством.
type FamilyRelation struct {
	FamilyID   int32  `json:"family_id"`
	FamilyName string `json:"family_name"`
	Score      int32  `json:"-"`
	//Explanation string `json:"explanation"`
}

// CropRelations – полный набор отношений для культуры.
type CropRelations struct {
	GoodPredecessors []CropRelation `json:"good_predecessors"`
	BadPredecessors  []CropRelation `json:"bad_predecessors"`
	GoodSuccessors   []CropRelation `json:"good_successors"`
	BadSuccessors    []CropRelation `json:"bad_successors"`
	GoodCompanions   []CropRelation `json:"good_companions"`
	BadCompanions    []CropRelation `json:"bad_companions"`

	GoodPredecessorFamilies []FamilyRelation `json:"good_predecessor_families"`
	BadPredecessorFamilies  []FamilyRelation `json:"bad_predecessor_families"`
	GoodSuccessorFamilies   []FamilyRelation `json:"good_successor_families"`
	BadSuccessorFamilies    []FamilyRelation `json:"bad_successor_families"`
	GoodCompanionFamilies   []FamilyRelation `json:"good_companion_families"`
	BadCompanionFamilies    []FamilyRelation `json:"bad_companion_families"`
}

// SunNeeds Enum: 1=shade, 2=partial, 3=full
type SunNeeds int32

const (
	SunNeedsShade   SunNeeds = 1
	SunNeedsPartial SunNeeds = 2
	SunNeedsFull    SunNeeds = 3
)

// ContextType Enum: 1=Predecessor, 2=Successor, 3=Companion
type ContextType int32

const (
	RuleContextPredecessor = 1 // Культура-предшественник (была до текущей)
	RuleContextSuccessor   = 2 // Культура-последователь (будет после текущей)
	RuleContextCompanion   = 3 // Культура-сосед
)
