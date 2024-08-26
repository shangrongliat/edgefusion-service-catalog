package config

type Config struct {
	Push Push `json:"push" yaml:"push"`
}

type Push struct {
	//云端直播服务地址
	CloudAddress string `json:"cloud_address" yaml:"cloud_address"`
	//云端是否存储
	IsCloudStorage bool `json:"is_cloud_storage" yaml:"is_cloud_storage"  default:"false"` //是否开启存储
	//云端直播是否开启
	IsCloudLive bool `json:"is_cloud_live" yaml:"is_cloud_live"  default:"false"` //是否开启直播
	//分发设置
	DistributionSetting bool `json:"distribution_setting" yaml:"distribution_setting" default:"false"` //是否开启分发
	//直播分发模式 直播代理 透传转发
	CloudLiveMode string `json:"cloud_live_mode" default:"0" yaml:"cloud_live_mode" default:"0"` //0 直播代理 1 透传转发
	With          int    `json:"with" yaml:"with" default:"1920"`                                //视频宽度
	Height        int    `json:"height" yaml:"height" default:"1080"`                            //视频高度
	Fps           int    `json:"fps" yaml:"fps" default:"30"`                                    //帧率
	InputSrc      string `json:"input_src" default:"" yaml:"input_src"`                          //输入源地址
}
