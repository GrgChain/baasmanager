package service

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"github.com/jonluo94/baasmanager/baas-gateway/model"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
)

type DashboardService struct {
	DbEngine *xorm.Engine
}

func NewDashboardService(engine *xorm.Engine) *DashboardService {
	return &DashboardService{
		DbEngine: engine,
	}
}

func (d *DashboardService) Counts(userAccount string) (bool, *model.Dashboard) {

	dash := new(model.Dashboard)
	var err error

	values := make([]interface{}, 0)
	where := "1=1"
	if userAccount != "" {
		where += " and user_account = ? "
		values = append(values, userAccount)
	}

	dash.Users, err = d.DbEngine.Count(new(entity.User))
	if err != nil {
		logger.Error(err.Error())
	}

	dash.Chains, err = d.DbEngine.Where(where, values...).Count(new(entity.Chain))
	if err != nil {
		logger.Error(err.Error())
	}
	dash.Chaincodes, err = d.DbEngine.Where(where, values...).Count(new(entity.Chaincode))
	if err != nil {
		logger.Error(err.Error())
	}
	dash.Channels, err = d.DbEngine.Where(where, values...).Count(new(entity.Channel))
	if err != nil {
		logger.Error(err.Error())
	}

	return true, dash
}

func (d *DashboardService) SevenDays(userAccount string, start, end int) (bool, map[string][]map[string]string) {

	sevenMap := make(map[string][]map[string]string)

	where := " where 1=1 "
	uwhere := where
	if userAccount != "" {
		where += fmt.Sprintf(" and user_account = '%s'" , userAccount)
	}

	if start != 0 {
		ws := fmt.Sprintf(" and created >= %d", start)
		where += ws
		uwhere += ws
	}

	if end != 0 {
		 ws := fmt.Sprintf(" and created <= %d", end)
		where += ws
		uwhere += ws
	}

	sql := ` SELECT from_unixtime( created, "%Y-%m-%d" ) AS days, count( id ) AS counts FROM `
	group := " GROUP BY days "
	table := "chain"
	chains, err := d.DbEngine.QueryString(sql + table + where + group)
	if err != nil {
		logger.Error(err.Error())
	}
	sevenMap["chains"] = chains

	table = "channel"
	channels, err := d.DbEngine.QueryString(sql + table + where + group)
	if err != nil {
		logger.Error(err.Error())
	}
	sevenMap["channels"] = channels

	table = "chaincode"
	chaincodes, err := d.DbEngine.QueryString(sql + table + where + group)
	if err != nil {
		logger.Error(err.Error())
	}
	sevenMap["chaincodes"] = chaincodes

	table = "user"
	users, err := d.DbEngine.QueryString(sql + table + uwhere + group)
	if err != nil {
		logger.Error(err.Error())
	}
	sevenMap["users"] = users

	return true, sevenMap
}

func (d *DashboardService) ConsensusTotal(userAccount string) (bool, []map[string]string) {

	sql := ` select count(1) as value ,consensus from chain `
	group := ` group by consensus `
	where := " where 1=1 "
	if userAccount != "" {
		where += fmt.Sprintf(" and user_account = '%s'" , userAccount)
	}

	totals, err := d.DbEngine.QueryString(sql + where + group)
	if err != nil {
		logger.Error(err.Error())
	}

	return true,totals
}