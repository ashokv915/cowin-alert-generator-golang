package models

type Centers struct {
	CenterID int64 `bson:"centerid"`
	Dose1 int64 `bson:"dose1"`
	Dose2 int64 `bson:"dose2"`
}