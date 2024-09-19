package consul

type Register struct {
	Datacenter     string     `json:"Datacenter"`
	ID             string     `json:"ID"`
	Node           string     `json:"Node"`
	Address        string     `json:"Address"`
	TaggedAddr     TaggedAddr `json:"TaggedAddresses"`
	NodeMeta       NodeMeta   `json:"NodeMeta"`
	Service        Service    `json:"Service"`
	Check          Check      `json:"Check"`
	SkipNodeUpdate bool       `json:"SkipNodeUpdate"`
}
type TaggedAddr struct {
	Lan string `json:"lan"`
	Wan string `json:"wan"`
}
type NodeMeta struct {
	Somekey string `json:"somekey"`
}
type Meta struct {
	RedisVersion string `json:"redis_version"`
}
type Service struct {
	ID         string       `json:"ID"`
	Service    string       `json:"Service"`
	Tags       []string     `json:"Tags"`
	Address    string       `json:"Address"`
	TaggedAddr TaggedAddr   `json:"TaggedAddresses"`
	Weights    AgentWeights `json:"Weights"`
	Meta       Meta         `json:"Meta"`
	Port       int          `json:"Port"`
	Namespace  string       `json:"Namespace"`
}
type Definition struct {
	TCP                            string `json:"TCP"`
	Interval                       string `json:"Interval"`
	Timeout                        string `json:"Timeout"`
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
}
type Check struct {
	Node       string     `json:"Node"`
	CheckID    string     `json:"CheckID"`
	Name       string     `json:"Name"`
	Notes      string     `json:"Notes"`
	Status     string     `json:"Status"`
	ServiceID  string     `json:"ServiceID"`
	Definition Definition `json:"Definition"`
	Namespace  string     `json:"Namespace"`
}

type AgentWeights struct {
	Passing int
	Warning int
}
