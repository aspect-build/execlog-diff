package types

type InputSet struct {
	Id              int64
	FileIds         []int64
	FileNames       []string
	NestedInputSets []int64
	LogPathHash     string
}
