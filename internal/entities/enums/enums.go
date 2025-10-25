package enums

import "time"

type ReversationStatus string

const (
	Booked    ReversationStatus = "booked"
	Cancelled ReversationStatus = "cancelled"
)

type Location string

const (
	Outside Location = "outside"
	Inside  Location = "inside"
)

type RevervationSlot int

const (
	Slot1 RevervationSlot = iota + 1
	Slot2
	Slot3
	Slot4
	Slot5
	Slot6
	Slot7
	Slot8
	Slot9
	Slot10
	Slot11
	Slot12
)

func GetSlot(datetime time.Time) RevervationSlot {

	hour := datetime.Hour()

	if hour >= 0 && hour < 2 {
		return Slot1
	} else if hour >= 2 && hour < 4 {
		return Slot2
	} else if hour >= 4 && hour < 6 {
		return Slot3
	} else if hour >= 6 && hour < 8 {
		return Slot4
	} else if hour >= 8 && hour < 10 {
		return Slot5
	} else if hour >= 10 && hour < 12 {
		return Slot6
	} else if hour >= 12 && hour < 14 {
		return Slot7
	} else if hour >= 14 && hour < 16 {
		return Slot8
	} else if hour >= 16 && hour < 18 {
		return Slot9
	} else if hour >= 18 && hour < 20 {
		return Slot10
	} else if hour >= 20 && hour < 22 {
		return Slot11
	} else if hour >= 22 && hour < 24 {
		return Slot12
	}
	return 0

}
