package app

type DBRole int16

const (
	ExistinDBRole DBRole = iota
	ErrorDBRole
	QueriedHistoryDBRole
)