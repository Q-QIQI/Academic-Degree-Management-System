package service

import (
	"data/baas-gateway/config"

	"github.com/jonluo94/baasmanager/baas-core/common/httputil"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
	"github.com/jonluo94/baasmanager/baas-core/core/model"
)

// logger：日志记录器，用于记录错误日志。日志级别设置为 ERROR，表示只记录错误信息。
var logger = log.GetLogger("service", log.ERROR)

// FabricService 是一个空结构体，用于组织与 Fabric 区块链引擎交互的方法。
type FabricService struct {
}

// 定义（创建）一条新的 Fabric 区块链。
// 调用 httputil.PostJson 方法，向 Fabric 引擎的 /defChain 接口发送 HTTP POST 请求。
// 请求的 URL 是通过 config.Config.GetString("BaasFabricEngine") 获取的 Fabric 引擎地址。
// 请求体是 chain 对象，会被序列化为JSON 格式。
func (g FabricService) DefChain(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/defChain", chain)
}

// 定义（创建）一个新的 Fabric 通道，并构建通道。
// 调用 httputil.PostJson 方法，向 Fabric 引擎的 /defChannelAndBuild 接口发送 HTTP POST 请求。
// 请求体是 chain 对象。
func (g FabricService) DefChannel(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/defChannelAndBuild", chain)
}

// 生成 Kubernetes 的 YAML 配置文件，并部署 Fabric 区块链。
// 调用 httputil.PostJson 方法，向 Fabric 引擎的 /defK8sYamlAndDeploy 接口发送 HTTP POST 请求。
// 请求体是 chain 对象
func (g FabricService) DeployK8sData(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/defK8sYamlAndDeploy", chain)
}

// 停止运行中的 Fabric 区块链。
func (g FabricService) StopChain(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/stopChain", chain)
}

// 释放（删除）Fabric 区块链。
func (g FabricService) ReleaseChain(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/releaseChain", chain)
}

func (g FabricService) DownloadChainArtifacts(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/downloadArtifacts", chain)
}

func (g FabricService) BuildChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/buildChaincode", channel)
}
func (g FabricService) UpdateChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/updateChaincode", channel)
}

// 查询链码信息。
func (g FabricService) QueryChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryChaincode", channel)
}

// 调用链码。
func (g FabricService) InvokeChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/invokeChaincode", channel)
}

func (g FabricService) UploadChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/uploadChaincode", channel)
}
func (g FabricService) DownloadChaincode(channel model.FabricChannel) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/downloadChaincode", channel)
}

func (g FabricService) QueryChainPods(chain model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryChainPods", chain)
}

func (g FabricService) QueryLedger(channel model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryLedger", channel)
}

func (g FabricService) QueryLatestBlocks(channel model.FabricChain) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryLatestBlocks", channel)
}

func (g FabricService) QueryBlock(channel model.FabricChain, search string) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/queryBlock?search="+search, channel)
}

// 修改 Fabric 区块链的 Kubernetes Pods 资源（如 CPU、内存等）。
func (g FabricService) ChangeChainPodResources(resource model.Resources) []byte {
	return httputil.PostJson(config.Config.GetString("BaasFabricEngine")+"/changeChainPodResources", resource)
}

func NewFabricService() *FabricService {
	return &FabricService{}
}
