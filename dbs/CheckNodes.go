package dbs

import (
	"check_vpn/dbs/model"
	"check_vpn/mylog"
	util "check_vpn/utils"
	"strings"
	"time"
)

func GetAllNode() (ret []*model.CheckNode, err error) {
	rd := GetDb()
	rq := rd.WithContext(util.GetCtx())
	ret, err = rq.CheckNode.Where(rd.CheckNode.IsShow.Eq(1)).Find()
	return
}
func GetOneNode(id int64) (ret *model.CheckNode, err error) {
	rd := GetDb()
	rq := rd.WithContext(util.GetCtx())
	ret, err = rq.CheckNode.Where(rd.CheckNode.IsShow.Eq(1)).Where(rd.CheckNode.ID.Eq(id)).First()
	return
}
func AddHist(newHis *model.CheckHistory) error {
	rd := GetDb()
	rq := rd.WithContext(util.GetCtx())

	return rq.CheckHistory.Create(newHis)
}
func TestAddhist(ip string) {
	//ip := "1.1.1.1"
	doid := int64(1)
	noid := int64(1)
	ret := util.Ping(ip)
	pinglong := float32(ret.AvgRtt.Seconds())

	newmod := model.CheckHistory{
		IP:       &ip,
		DonodeID: &doid,
		NodeID:   &noid,
		PingLong: &pinglong,
		NodeLong: &pinglong,
		APILong:  &pinglong,
	}

	mylog.Der(AddHist(&newmod))
	mylog.Logf("id:%v", newmod.ID)
}
func DoCheckWg(body util.WgItem, one model.CheckNode, api_long float32) {
	ips := strings.Split(body.UserEndpoint, ":")
	pingret := util.Ping(ips[0])
	pinglong := float32(pingret.AvgRtt.Seconds())
	nodlong := float32(0)
	if pinglong != 0 {
		nodlong = util.CheckWg(body.UserAddress, body.UserPrivateKey, body.ServerPublicKey, body.AllowedIPs, body.UserEndpoint)
	}
	newmod := model.CheckHistory{
		IP:       &ips[0],
		DonodeID: &util.DoId,
		NodeID:   &one.ID,
		PingLong: &pinglong,
		NodeLong: &nodlong,
		APILong:  &api_long,
	}
	mylog.Der(AddHist(&newmod))
	mylog.Logf("id:%v", newmod.ID)
}
func DoCheckSs(body util.SsItem, one model.CheckNode, api_long float32) {
	pingret := util.Ping(body.Ip)
	pinglong := float32(pingret.AvgRtt.Seconds())
	nodlong := float32(0)
	if pinglong != 0 {
		nodlong = util.CheckSs(body.Ip+":"+body.Port, body.Method, body.Paskey)

	}
	newmod := model.CheckHistory{
		IP:       &body.Ip,
		DonodeID: &util.DoId,
		NodeID:   &one.ID,
		PingLong: &pinglong,
		NodeLong: &nodlong,
		APILong:  &api_long,
	}
	mylog.Der(AddHist(&newmod))
	mylog.Logf("id:%v", newmod.ID)
}
func DoCheck(one model.CheckNode) {
	url := *one.Host + *one.ListPath
	startt := time.Now().UnixMilli()
	if *one.NodeType == "wg" {
		mylog.Logf("GetWg:url:%v", url)

		body, err := util.GetWgRet(url, *one.ReqEncode, *one.ReqEncodeKey)
		api_long := float32(time.Now().UnixMilli()-startt) / float32(1000.0000000000000)

		if len(body.Test) == 0 {
			mylog.Logf("GetWgRet:err:%v", err)
		}
		for _, itm := range body.Test {
			mylog.Logf("GetWgRet:item:%v", itm)
			DoCheckWg(itm, one, api_long)
		}
	}
	if *one.NodeType == "ss" {
		body, err := util.GetSsRet(url, *one.ReqEncode, *one.ReqEncodeKey)
		api_long := float32(time.Now().UnixMilli()-startt) / float32(1000.0000000000000)

		if len(body.Test) == 0 {
			mylog.Logf("GetWgRet:err:%v", err)
		}
		for _, itm := range body.Test {
			mylog.Logf("GetWgRet:item:%v", itm)
			DoCheckSs(itm, one, api_long)
		}
	}
}
