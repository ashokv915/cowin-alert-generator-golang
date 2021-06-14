package models

type JSONResponse struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	Center_id int64 `json:"center_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	StateName string `json:"state_name"`
	DistrictName string `json:"district_name"`
	BlockName string `json:"block_name"`
	Pincode int64 `json:"pincode"`
	From string `json:"from"`
	To string `json:"to"`
	Lat int64 `json:"lat"`
	Long int64 `json:"long"`
	FeeType string `json:"fee_type"`
	SessionID string `json:"session_id"`
	Date string `json:"date"`
	AvailableCapacity int16 `json:"available_capacity"`
	AvailableCapacityDose1 int16 `json:"available_capacity_dose1"`
	AvailableCapacityDose2 int16 `json:"available_capacity_dose2"`
	Fee string `json:"fee"`
	MinAge int16 `json:"min_age_limit"`
	Vaccine string `json:"vaccine"`
	Slots []string `json:"slots"`
}