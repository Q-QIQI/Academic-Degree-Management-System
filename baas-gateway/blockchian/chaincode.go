// InstallCC
// InstantiateCC
// UpgradeCC
// InvokeCC
// InvokeCCDelete
// QueryCC
// QueryCCInfo
// GetBlocks
// Close
// packArgs
package blockchian

import (
	"log"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric-protos-go/common"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// InstallCC install chaincode for target peer
func (c *Client) InstallCC(v string, peer string) error {
	targetPeer := resmgmt.WithTargetEndpoints(peer)

	// pack the chaincode
	ccPkg, err := gopackager.NewCCPackage(c.CCPath, c.CCGoPath)
	if err != nil {
		return errors.WithMessage(err, "pack chaincode error")
	}

	// new request of installing chaincode
	req := resmgmt.InstallCCRequest{
		Name:    c.CCID,   //智能合约ID
		Path:    c.CCPath, //智能合约路径
		Version: v,        //版本号
		Package: ccPkg,    //打包后的智能合约
	}

	//发送安装请求
	resps, err := c.resMgtClient.InstallCC(req, targetPeer)
	if err != nil {
		return errors.WithMessage(err, "installCC error")
	}

	// 处理安装响应
	//检查 resp.Status 是否为 http.StatusOK，否则记录错误。
	//如果 resp.Info 为 "already installed"，说明该智能合约已安装，直接返回 nil。
	var errs []error
	for _, resp := range resps {
		log.Printf("Install  response status: %v", resp.Status)
		if resp.Status != http.StatusOK {
			errs = append(errs, errors.New(resp.Info))
		}
		if resp.Info == "already installed" {
			log.Printf("Chaincode %s already installed on peer: %s.\n",
				c.CCID+"-"+v, resp.Target)
			return nil
		}
	}

	if len(errs) > 0 {
		log.Printf("InstallCC errors: %v", errs)
		return errors.WithMessage(errs[0], "installCC first error")
	}
	return nil
}

// 实例化智能合约
func (c *Client) InstantiateCC(v string, peer string) (fab.TransactionID,
	error,
) {
	// endorser policy
	//设置背书策略，设置组织1或组织2都能背书
	org1OrOrg2 := "OR('org1MSP.member','org2MSP.member')"
	// org1OrOrg2 := "OR('Org1MSP.member','Org2MSP.member')"
	//C.genPolicy解析该策略
	ccPolicy, err := c.genPolicy(org1OrOrg2)
	if err != nil {
		return "", errors.WithMessage(err, "gen policy from string error")
	}

	// new request样例请求
	// Attention: args should include `init` for Request not
	// have a method term to call init
	args := packArgs([]string{"john", "1"}) //链码初始化
	req := resmgmt.InstantiateCCRequest{
		Name:    c.CCID,
		Path:    c.CCPath,
		Version: v,
		Args:    args,
		Policy:  ccPolicy,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.resMgtClient.InstantiateCC(c.ChannelID, req, reqPeers)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return "", nil
		}
		return "", errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return resp.TransactionID, nil //如果链码已实例化，返回nil，否则返回Transaction
}

func (c *Client) genPolicy(p string) (*common.SignaturePolicyEnvelope, error) {
	// TODO bug, this any leads to endorser invalid
	if p == "ANY" {
		return cauthdsl.SignedByAnyMember([]string{c.OrgName}), nil
	}
	return cauthdsl.FromString(p)
}

// 调用智能合约
func (c *Client) InvokeCC(peers []string, funcName string, params []string) (fab.TransactionID, error) {
	// new channel request for invoke
	args := packArgs(params)
	req := channel.Request{
		ChaincodeID: c.CCID,
		Fcn:         funcName,
		Args:        args,
	}

	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.channelClient.Execute(req, reqPeers)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) InvokeCCDelete(peers []string) (fab.TransactionID, error) {
	log.Println("Invoke delete")
	// new channel request for invoke
	args := packArgs([]string{"c"})
	req := channel.Request{
		ChaincodeID: c.CCID,
		Fcn:         "delete",
		Args:        args,
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.channelClient.Execute(req, reqPeers)
	log.Printf("Invoke chaincode delete response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)
	if err != nil {
		return "", errors.WithMessage(err, "invoke chaincode error")
	}

	return resp.TransactionID, nil
}

func (c *Client) QueryCC(peer, funcName, key string) ([]byte, error) {
	// new channel request for query
	req := channel.Request{
		ChaincodeID: c.CCID,
		Fcn:         funcName,
		Args:        packArgs([]string{key}),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer)
	resp, err := c.channelClient.Query(req, reqPeers)
	if err != nil {
		return nil, errors.WithMessage(err, "query chaincode error")
	}

	return resp.Payload, nil
}

func (c *Client) UpgradeCC(v string, peer string) error {
	// endorser policy
	org1AndOrg2 := "AND('Org1MSP.member','Org2MSP.member')"
	ccPolicy, err := c.genPolicy(org1AndOrg2)
	if err != nil {
		return errors.WithMessage(err, "gen policy from string error")
	}

	// new request
	// Attention: args should include `init` for Request not
	// have a method term to call init
	// Reset a b's value to test the upgrade
	args := packArgs([]string{"set", "a"})
	req := resmgmt.UpgradeCCRequest{
		Name:    c.CCID,
		Path:    c.CCPath,
		Version: v,
		Args:    args,
		Policy:  ccPolicy,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.resMgtClient.UpgradeCC(c.ChannelID, req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return nil
}

func (c *Client) QueryCCInfo(v string, peer string) {
}

func (c *Client) GetBlocks() {
}

func (c *Client) Close() {
	c.SDK.Close()
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
