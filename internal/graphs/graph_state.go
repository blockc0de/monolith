package graphs

type GraphStateEnum int

const (
	GraphStateEnumStopped    GraphStateEnum = 0
	GraphStateEnumStarting   GraphStateEnum = 1
	GraphStateEnumStarted    GraphStateEnum = 2
	GraphStateEnumError      GraphStateEnum = 3
	GraphStateEnumRestarting GraphStateEnum = 4
	GraphStateEnumPause      GraphStateEnum = 5
)
