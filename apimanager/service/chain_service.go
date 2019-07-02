package service

import (
	"gitee.com/jonluo/baasmanager/apimanager/entity"
	"github.com/go-xorm/xorm"
	"github.com/jonluo94/commontools/gintool"
	"encoding/json"
	"gitee.com/jonluo/baasmanager/apimanager/model"
	"bytes"
	"io"
)

type ChainService struct {
	DbEngine      *xorm.Engine
	FabircService *FabricService
}

func (l *ChainService) Add(chain *entity.Chain) (bool, string) {

	i, err := l.DbEngine.Insert(chain)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *ChainService) Update(chain *entity.Chain) (bool, string) {

	i, err := l.DbEngine.Where("id = ?", chain.Id).Update(chain)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *ChainService) UpdateStatus(chain *entity.Chain) (bool, string) {

	sql := "update `chain` set status = ? where id = ?"
	res, err := l.DbEngine.Exec(sql, chain.Status, chain.Id)
	if err != nil {
		logger.Error(err.Error())
	}
	r,err := res.RowsAffected()
	if err == nil && r > 0{
		return true, "update success"
	}
	return false, "update fail"
}


func (l *ChainService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Chain{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

func (l *ChainService) GetByChain(chain *entity.Chain) (bool, *entity.Chain) {
	has, err := l.DbEngine.Get(chain)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, chain
}

func (l *ChainService) GetList(chain *entity.Chain, page, size int) (bool, []*entity.Chain, int64) {
	pager := gintool.CreatePager(page, size)
	chains := make([]*entity.Chain, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if chain.Name != "" {
		where += " and name = ? "
		values = append(values, chain.Name)
	}
	if chain.UserAccount != "" {

		where += " and user_account = ? "
		values = append(values, chain.UserAccount)
	}
	if chain.Consensus != "" {
		where += " and consensus = ? "
		values = append(values, chain.Consensus)
	}
	if chain.PeersOrgs != "" {
		where += " and peers_orgs like ? "
		values = append(values, "%"+chain.PeersOrgs+"%")
	}
	if chain.TlsEnabled != "" {
		where += " and tls_enabled = ? "
		values = append(values, chain.TlsEnabled)
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&chains)
	if err != nil {
		logger.Error(err.Error())
	}
	total, err := l.DbEngine.Where(where, values...).Count(new(entity.Chain))
	if err != nil {
		logger.Error(err.Error())
	}
	return true, chains, total
}

func (l *ChainService) BuildChain(chain *entity.Chain) (bool, string) {

	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.DefChain(fc)
	var ret model.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "build fail"
	}

	if ret.Code == 0 {
		chain.Status = 1
		return l.UpdateStatus(chain)
	} else {
		return false, "build fail"
	}

}

func (l *ChainService) RunChain(chain *entity.Chain) (bool, string) {

	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.DeployK8sData(fc)
	var ret model.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "run fail"
	}

	if ret.Code == 0 {
		chain.Status = 2
		return l.UpdateStatus(chain)
	} else {
		return false, "build fail"
	}

}

func (l *ChainService) StopChain(chain *entity.Chain) (bool, string) {

	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.StopChain(fc)
	var ret model.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "run fail"
	}

	if ret.Code == 0 {
		chain.Status = 3
		return l.UpdateStatus(chain)
	} else {
		return false, "build fail"
	}

}

func (l *ChainService) ReleaseChain(chain *entity.Chain) (bool, string) {

	fc := entity.ParseFabircChain(chain)
	resp := l.FabircService.ReleaseChain(fc)
	var ret model.RespData
	err := json.Unmarshal(resp, &ret)
	if err != nil {
		return false, "run fail"
	}

	if ret.Code == 0 {
		chain.Status = 0
		return l.UpdateStatus(chain)
	} else {
		return false, "build fail"
	}

}

func (l *ChainService) DownloadChainArtifacts(chain *entity.Chain) (io.Reader, int64, string) {

	fc := entity.ParseFabircChain(chain)
	bts := l.FabircService.DownloadChainArtifacts(fc)
	reader := bytes.NewReader(bts)
	contentLength := reader.Len()
	return reader, int64(contentLength), chain.Name + ".tar"

}

func NewChainService(engine *xorm.Engine, fabircService *FabricService) *ChainService {
	return &ChainService{
		DbEngine:      engine,
		FabircService: fabircService,
	}
}
