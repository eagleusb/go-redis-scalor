package rediscluster

type RedisNode struct {
	Slots           []string
	SlotsCount      int
	SlotsPercentage int
	Id              string
	Ip              string
	Name            string
}

type RedisNodes struct {
	ClusterNodes map[string]RedisNode
}

type RedisClusterConf struct {
	Conf          map[string]string
	State         string
	Size          string
	SlotsAssigned string
	SlotsOk       string
	NodesKnown    string
}
