package model

type SlotId string
type BannerId string
type GroupId string

type StatTemp struct {
	Id            BannerId
	CountDisplays int64
	CountClicks   int64
}
