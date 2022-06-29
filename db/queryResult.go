package db

type DBQueryResult struct {
	DBRole      DBRole `json:"db_role"`
	QueryString string `json:"query_string"`
	QueryResult int    `json:"query_result"`
}
