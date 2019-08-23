package service

import (
	"github.com/go-xorm/xorm"
	"time"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-core/common/json"
)

type ChaincodeService struct {
	DbEngine      *xorm.Engine
	FabircService *FabricService
}

func (l *ChaincodeService) Add(cc *entity.Chaincode) (bool, string) {
	cc.Created = time.Now().Unix()
	cc.Version = "1"
	cc.Status = 0
	i, err := l.DbEngine.Insert(cc)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *ChaincodeService) Update(cc *entity.Chaincode) (bool, string) {

	i, err := l.DbEngine.Where("id = ?", cc.Id).Update(cc)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *ChaincodeService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Chaincode{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

func (l *ChaincodeService) GetByChaincode(cc *entity.Chaincode) (bool, *entity.Chaincode) {
	has, err := l.DbEngine.Get(cc)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, cc
}

func (l *ChaincodeService) GetList(cc *entity.Chaincode, page, size int) (bool, []*entity.Chaincode, int64) {
	pager := gintool.CreatePager(page, size)
	ccs := make([]*entity.Chaincode, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if cc.ChaincodeName != "" {
		where += " and chaincode_name = ? "
		values = append(values, cc.ChaincodeName)
	}
	if cc.ChannelId != 0 {
		where += " and channel_id = ? "
		values = append(values, cc.ChannelId)
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&ccs)
	if err != nil {
		logger.Error(err.Error())
	}
	total, err := l.DbEngine.Where(where, values...).Count(new(entity.Chaincode))
	if err != nil {
		logger.Error(err.Error())
	}
	return true, ccs, total
}

func (l *ChaincodeService) GetAllList(chainId int) (bool, []*entity.Chaincode) {

	ccs := make([]*entity.Chaincode, 0)
	err := l.DbEngine.Where("channel_id = ?", chainId).Find(&ccs)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, ccs
}

func (l *ChaincodeService) AddChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {

	bys, err := ioutil.ReadFile(cc.GithubPath)
	if err != nil {
		return false, "add fail"
	}
	cc.Version = "1"
	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	fc.ChaincodeBytes = bys
	resp := l.FabircService.UploadChaincode(fc)
	var ret gintool.RespData
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "add fail"
	}

	if ret.Code == 0 {
		cc.GithubPath = ret.Data.(string)
		cc.Created = time.Now().Unix()
		cc.Status = 0
		return l.Add(cc)
	} else {
		return false, "add fail"
	}
}

func (l *ChaincodeService) DeployChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {

	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	args := make([][]byte, 1)
	args[0] = []byte("init")

	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp := l.FabircService.BuildChaincode(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "deploy fail"
	}

	if ret.Code == 0 {
		cc.Status = 1
		return l.Update(cc)
	} else {
		return false, "deploy fail"
	}
}

func (l *ChaincodeService) UpgradeChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {

	bys, err := ioutil.ReadFile(cc.GithubPath)
	if err != nil {
		return false, "upgrade fail"
	}
	v, err := strconv.Atoi(cc.Version)
	if err != nil {
		return false, "version error"
	}
	cc.Version = fmt.Sprintf("%d", v+1)

	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	fc.ChaincodeBytes = bys
	resp := l.FabircService.UploadChaincode(fc)
	var ret gintool.RespData
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "upgrade fail"
	}

	if ret.Code == 0 {
		cc.GithubPath = ret.Data.(string)
		fc.ChaincodePath = ret.Data.(string)
	} else {
		return false, "upload fail"
	}

	args := make([][]byte, 1)
	args[0] = []byte("init")
	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp = l.FabircService.UpdateChaincode(fc)
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "upgrade fail"
	}

	if ret.Code == 0 {
		return l.Update(cc)
	} else {
		return false, "upgrade fail"
	}
}

func (l *ChaincodeService) InvokeChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {

	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	args := make([][]byte, 0)

	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp := l.FabircService.InvokeChaincode(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "invoke fail"
	}

	if ret.Code == 0 {
		return true, ret.Data.(string)
	} else {
		return false, "invoke fail"
	}
}

func (l *ChaincodeService) QueryChaincode(chain *entity.Chain, channel *entity.Channel, cc *entity.Chaincode) (bool, string) {

	fc := entity.ParseFabircChannel(entity.ParseFabircChainAndChannel(chain, channel), cc)
	args := make([][]byte, 0)

	for _, v := range strings.Split(cc.Args, ",") {
		args = append(args, []byte(v))
	}
	fc.Args = args
	resp := l.FabircService.QueryChaincode(fc)
	var ret gintool.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "query fail"
	}

	if ret.Code == 0 {
		return true, ret.Data.(string)
	} else {
		return false, "query fail"
	}
}

func NewChaincodeService(engine *xorm.Engine, fabircService *FabricService) *ChaincodeService {
	return &ChaincodeService{
		DbEngine:      engine,
		FabircService: fabircService,
	}
}
