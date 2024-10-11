package models

type EngineType int8

const (
	HTTP EngineType = iota
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
