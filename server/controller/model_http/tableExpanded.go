package model_http

type TableState string

const (
	TableStateExpanded TableState = "expanded"
	TableStateNormal   TableState = ""
	TableStateDelta    TableState = "delta"
)

type TableMode struct {
	TableState TableState `cookie:"expand-state"`
}

type TableSetMode struct {
	ExpandTable TableState `form:"expand-table"`
}
