package clusterState

type ClusterState int

const (
	Follow ClusterState = 2
	Leader ClusterState = 4
)
