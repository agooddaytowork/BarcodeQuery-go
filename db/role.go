package db

type DBRole int16

const (
	ExistingDBRole DBRole = iota
	ErrorDBRole
	QueriedHistoryDBRole
)