package consul

type Register struct {
	ID                string            `json:"ID"`
	Name              string            `json:"Name"`
	Tags              *[]string         `json:"Tags"`
	Address           string            `json:"Address"`
	Port              int               `json:"Port"`
	Meta              map[string]string `json:"Meta"`
	Check             *Check            `json:"Check"`
	EnableTagOverride bool              `json:"EnableTagOverride"`
	Weights           Weights           `json:"Weights"`
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"`
	Args                           []string `json:"Args"`
	Interval                       string   `json:"Interval"`
	Timeout                        string   `json:"Timeout"`
}

type Weights struct {
	Passing int `json:"Passing"`
	Warning int `json:"Warning"`
}
