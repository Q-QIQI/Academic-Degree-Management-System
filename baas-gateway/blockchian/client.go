// 安装，调用，查询
package blockchian

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// 定义Fabric的客户端信息
type Client struct {
	// Fabric network information
	ConfigPath string // Fabric 配置文件路径
	OrgName    string // 组织名称
	OrgAdmin   string // 组织管理员用户名
	OrgUser    string // 组织普通用户

	// sdk clients
	SDK           *fabsdk.FabricSDK // Hyperledger Fabric SDK 实例
	resMgtClient  *resmgmt.Client   // 资源管理客户端（用于管理通道、安装链码等）
	channelClient *channel.Client   // 通道客户端（用于与链码交互）

	// Same for each peer
	ChannelID string
	CCID      string // chaincode ID, eq name
	CCPath    string // chaincode source path, 是GOPATH下的某个目录
	CCGoPath  string // GOPATH used for chaincode
}

// 初始化 Fabric 客户端
func New(cfg, org, admin, user string) *Client {
	// 初始化 Client 结构体
	c := &Client{
		ConfigPath: cfg, //指定 Fabric 网络的 YAML 配置文件路径
		OrgName:    org,
		OrgAdmin:   admin,
		OrgUser:    user,

		CCID:     "sacc",
		CCPath:   "chaincode/sacc/", // 相对路径是从GOPAHT/src开始的
		CCGoPath: os.Getenv("GOPATH"),
		// ChannelID: "application-channel-1",
		ChannelID: "mychannelone",
	}

	// create sdk
	//通过 Fabric SDK 加载 ConfigPath 指定的配置文件，创建 SDK 实例。
	//fabsdk.New() 会解析 YAML 文件，连接到 Fabric 网络。
	sdk, err := fabsdk.New(config.FromFile(c.ConfigPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}
	c.SDK = sdk
	log.Println("Initialized fabric sdk")

	//创建资源管理客户端
	rcp := sdk.Context(fabsdk.WithUser(c.OrgAdmin), fabsdk.WithOrg(c.OrgName))
	c.resMgtClient, err = resmgmt.New(rcp)
	if err != nil {
		log.Panicf("failed to create resource client: %s", err)
	}
	log.Println("Initialized resource client")
	ccp := sdk.ChannelContext(c.ChannelID, fabsdk.WithUser(c.OrgUser), fabsdk.WithOrg(c.OrgName))
	c.channelClient, err = channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}
	log.Println("Initialized channel client")

	return c
}
