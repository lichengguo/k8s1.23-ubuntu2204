package nacoscfg

import (
	"go-k8s-one/conf"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func LoadNacos() []byte {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(conf.IP, conf.Port, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(conf.NameSpaceID),
		//constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		//constant.WithLogDir("/tmp/nacos/log"),
		//constant.WithCacheDir("/tmp/nacos/cache"),
		//constant.WithLogLevel("error"),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: conf.DataID,
		Group:  conf.Group,
	})
	if err != nil {
		panic(err)
	}

	return []byte(content)
}
