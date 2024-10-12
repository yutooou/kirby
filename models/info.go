package models

type EngineType int8

const (
	HTTP EngineType = iota
)

var (
	KirbyInfo Info = Info{
		Name:       "Kirby",
		Code:       "_kirby",
		Version:    "v0.0.1",
		Desc:       "Kirby, data simulation engine.",
		EngineType: HTTP,
	}
)

type (
	// control point static infomation
	Info struct {
		// name
		Name string
		// unique identifier of the control point
		Code string
		// version
		Version string
		// description
		Desc string
		// engine type
		EngineType EngineType
	}
)
