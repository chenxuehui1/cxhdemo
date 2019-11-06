package service

import (
	"encoding/json"
	"os"

	"github.com/chenxuehui1/cxhdemo/common"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hashicorp/memberlist"
)

type delegateCxhdemo struct {
	broadcasts *memberlist.TransmitLimitedQueue
	extendIP   string
	logger     log.Logger
}

//deleageteWatchmen工厂
func newDelegateCxhdemo(m *memberlist.Memberlist) *delegateCxhdemo {
	d := &delegateCxhdemo{}
	d.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	d.logger =log.NewLogfmtLogger(os.Stdout)
	//d.logger = level.NewFilter(d.logger, Opt)

	return d

}

//memberlist delegate接口必须实现的方法(1/5)
func (d *delegateCxhdemo) NodeMeta(limit int) []byte {
	tmp := common.MateData{
		ExtendIP: "",
		Label:    "watchmen",
	}
	data, _ := json.Marshal(&tmp)
	return data
}

//memberlist delegate接口必须实现的方法(2/5)
func (d *delegateCxhdemo) NotifyMsg(b []byte) {
}

//memberlist delegate接口必须实现的方法(3/5)
func (d *delegateCxhdemo) GetBroadcasts(overhead, limit int) [][]byte {
	return d.broadcasts.GetBroadcasts(overhead, limit)
}

//memberlist delegate接口必须实现的方法(4/5)
//join为true是向所有节点push/pull，为false时定期随机选择一个节点push/pull
func (d *delegateCxhdemo) LocalState(join bool) []byte {
	return []byte(d.extendIP)
}

//memberlist delegate接口必须实现的方法(5/5)
//接收到对端的LocalState消息会触发本方法
func (d *delegateCxhdemo) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	d.extendIP = string(buf)
	level.Debug(d.logger).Log("methodName", "MergeRemoteState", "WatchmenGetExtendIP", d.extendIP)
}
