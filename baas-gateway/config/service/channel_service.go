package service

import (
	"data/baas-gateway/entity"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
)

type ChannelService struct {
	DbEngine      *xorm.Engine   //xorm的数据库引擎，用于执行数据库操作。
	FabircService *FabricService //xorm的数据库引擎，用于执行数据库操作。
}

// 向数据库中插入一条新的Channel记录。
func (l *ChannelService) Add(channel *entity.Channel) (bool, string) {

	i, err := l.DbEngine.Insert(channel)

	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

// 更新数据库中已有的Channel记录。
func (l *ChannelService) Update(channel *entity.Channel) (bool, string) {

	i, err := l.DbEngine.Where("id = ?", channel.Id).Update(channel)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

// 根据id删除数据库中的Channel记录。
func (l *ChannelService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Channel{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

// 根据chain_id删除与该链相关的所有Channel记录，并且删除与之关联的chaincode记录。
func (l *ChannelService) DeleteByChainId(id int) (bool, string) {

	sql := "delete from chaincode where channel_id in ( select id from channel where chain_id = ?)"
	_, err := l.DbEngine.Exec(sql, id)
	if err != nil {
		logger.Error(err.Error())
	}

	i, err := l.DbEngine.Where("chain_id = ?", id).Delete(&entity.Channel{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

// 根据传入的channel对象查询数据库中的记录。
func (l *ChannelService) GetByChannel(channel *entity.Channel) (bool, *entity.Channel) {
	fmt.Printf("GetByChannel channel channel.Id = <%d> channel.ChainId = <%d> channel.Orgs = <%s> channel.ChannelName = <%s> channel.UserAccount = <%s> channel.Created = <%d>", channel.Id, channel.ChainId, channel.Orgs, channel.ChannelName, channel.UserAccount, channel.Created)
	has, err := l.DbEngine.Get(channel)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, channel
}

// 根据channel的Id查询数据库中的记录。
func (l *ChannelService) GetByChannelId(channel *entity.Channel) (bool, *entity.Channel) {
	fmt.Printf("GetByChannelId channel channel.Id = <%d>\n", channel.Id)
	var ch entity.Channel
	_, err := l.DbEngine.Id(channel.Id).Get(&ch)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Printf("channel Id = <%d> ChainId = <%d> Orgs = <%s> ChannelName = <%s> UserAccount = <%s> Created = <%d>\n", ch.Id, ch.ChainId, ch.Orgs, ch.ChannelName, ch.UserAccount, ch.Created)
	return true, &ch
}

// 根据分页参数查询Channel列表。
func (l *ChannelService) GetList(channel *entity.Channel, page, size int) (bool, []*entity.Channel) {

	channels := make([]*entity.Channel, 0)

	values := make([]interface{}, 0)

	where := "1=1"

	err := l.DbEngine.Where(where, values...).Limit(size, page).Find(&channels)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, channels
}

// 根据chainId查询所有相关的Channel记录。
func (l *ChannelService) GetAllList(chainId int) (bool, []*entity.Channel) {

	channels := make([]*entity.Channel, 0)
	err := l.DbEngine.Where("chain_id = ?", chainId).Find(&channels)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, channels
}

// 添加一个新的Channel记录，并且调用FabricService来定义Fabric区块链中的通道。
func (l *ChannelService) AddChannel(chain *entity.Chain, channel *entity.Channel) (bool, string) {

	fc := entity.ParseFabircChainAndChannel(chain, channel)
	resp := l.FabircService.DefChannel(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "add fail"
	}

	if ret.Code == 0 {
		channel.Created = time.Now().Unix()
		return l.Add(channel)
	} else {
		return false, "add fail"
	}

}

// 创建一个新的ChannelService实例。
func NewChannelService(engine *xorm.Engine, fabircService *FabricService) *ChannelService {
	return &ChannelService{
		DbEngine:      engine,
		FabircService: fabircService,
	}
}
