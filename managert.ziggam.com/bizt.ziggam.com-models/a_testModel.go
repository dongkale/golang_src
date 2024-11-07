package models

// RoomTest
type RoomTest struct {
	MemNo  string
	MemId  string
	MemSn  string
	AuthCd int
}

var Room1 RoomTest = RoomTest{"a", "b", "c", 33}

// TestID
var TestID = map[string]string{
	"Created":  "Created",
	"Running":  "Running",
	"Finished": "Finished",
	"Errorred": "Errorred",
}
