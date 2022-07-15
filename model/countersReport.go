package model

type CounterReport struct {
	QueryCounter             int `json:"query_counter"`
	QueryCounterLimit        int `json:"query_counter_limit"`
	PackageCounter           int `json:"package_counter"`
	TotalCounter             int `json:"total_counter"`
	NumberOfItemInExistingDB int `json:"number_of_item_in_existing_db"`
	NumberOfCameraScanError  int `json:"number_of_camera_scan_error"`
}
