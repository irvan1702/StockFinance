package entity

import "time"

type Stock struct {
	Id      int
	Date    time.Time
	Open    int
	High    int
	Low     int
	Close   int
	Volume  int
	Action  string
	Summary int
}
