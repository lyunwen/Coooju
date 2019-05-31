package global

import (
	"../models"
	"../models/clusterState"
)

//全局异常
var (
	Errors      []error
	ClusterData *models.Data
	CurrentData *CurrentNodeInfo
)

type CurrentNodeInfo struct {
	VotedTerm    int
	VotedState   VotedState
	ClusterState clusterState.ClusterState
	Name         string
	Term         int
	Address      string
}

type VotedState int

const (
	VotedState_UnDo VotedState = 1
	VotedState_Done VotedState = 2
)
