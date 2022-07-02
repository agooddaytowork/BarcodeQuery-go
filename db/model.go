package db

type QueryResult struct {
	DBRole      DBRole `json:"db_role"`
	QueryString string `json:"query_string"`
	QueryResult int    `json:"query_result"`
}

type StateUpdate struct {
	DBRole DBRole        `json:"db_role"`
	State  []QueryResult `json:"state"`
}
