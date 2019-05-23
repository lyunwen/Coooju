package SelfFlagStatus

type SelfFlagStatus int

const (
	Init      SelfFlagStatus = 1
	Follow    SelfFlagStatus = 2
	Candidate SelfFlagStatus = 3
	Leader    SelfFlagStatus = 4
)
