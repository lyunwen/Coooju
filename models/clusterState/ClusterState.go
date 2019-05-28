package clusterState

type ClusterState int32

const (
	Follow    ClusterState = 1
	Candidate ClusterState = 2
	Leader    ClusterState = 4
)
