package service

import (
	"gitee.com/jonluo/baasmanager/apimanager/entity"
	"github.com/go-xorm/xorm"
	"gitee.com/jonluo/baasmanager/apimanager/model"
	"encoding/json"
	"time"
)

type ChannelService struct {
	DbEngine      *xorm.Engine
	FabircService *FabricService
}

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

func (l *ChannelService) GetByChannel(channel *entity.Channel) (bool, *entity.Channel) {
	has, err := l.DbEngine.Get(channel)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, channel
}

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

func (l *ChannelService) GetAllList(chainId int) (bool, []*entity.Channel) {

	channels := make([]*entity.Channel, 0)
	err := l.DbEngine.Where("chain_id = ?", chainId).Find(&channels)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, channels
}

func (l *ChannelService) AddChannel(chain *entity.Chain, channel *entity.Channel) (bool, string) {

	fc := entity.ParseFabircChainAndChannel(chain, channel)
	resp := l.FabircService.DefChannel(fc)
	var ret model.RespData
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

func NewChannelService(engine *xorm.Engine, fabircService *FabricService) *ChannelService {
	return &ChannelService{
		DbEngine:      engine,
		FabircService: fabircService,
	}
}
