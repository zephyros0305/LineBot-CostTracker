package main

type OperationType int

const (
	Error OperationType = iota
	KeepRecord
	GetRecord
	GetStatistic
)

type OperationData struct {
	CostType string
	Number   uint64
	Memo     string
}
