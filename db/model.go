package db

type QueryIntResult struct {
	DBRole      DBRole `json:"db_role"`
	QueryString string `json:"query_string"`
	QueryResult int    `json:"query_result"`
}

type QueryStringResult struct {
	DBRole      DBRole `json:"db_role"`
	QueryString string `json:"query_string"`
	QueryResult string `json:"query_result"`
}

type StateUpdate struct {
	DBRole DBRole           `json:"db_role"`
	State  []QueryIntResult `json:"state"`
}
