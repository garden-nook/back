package enum

// RuleContextType
// @Description Enum: RuleContextPredecessor=1, RuleContextCompanion=2
type RuleContextType int32

const (
	RuleContextPredecessor RuleContextType = 1 // Культура-предшественник (была до текущей)
	RuleContextCompanion   RuleContextType = 2 // Культура-сосед
)
