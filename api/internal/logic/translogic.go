package logic

import (
	"context"
	"fmt"
	"github.com/dtm-labs/dtmdriver"
	driver "github.com/dtm-labs/dtmdriver-gozero"
	"github.com/dtm-labs/dtmgrpc"
	"tcc/api/internal/svc"
	"tcc/api/internal/types"
	"tcc/trans/transclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransLogic(ctx context.Context, svcCtx *svc.ServiceContext) TransLogic {
	return TransLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransLogic) Trans(req types.TransRequest) (resp *types.TransResponse, err error) {
	transRpcBusiServer, err := l.svcCtx.Config.TransRpc.BuildTarget()
	if err != nil {
		return nil, err
	}
	// dtm 服务的 etcd 注册地址
	var dtmServer = "etcd://localhost:2379/dtmservice"
	// 创建一个gid
	gid := dtmgrpc.MustGenGid(dtmServer)
	if err := dtmdriver.Use(driver.DriverName); err != nil {
		return nil, err
	}
	err = dtmgrpc.TccGlobalTransaction(dtmServer, gid, func(tcc *dtmgrpc.TccGrpc) error {
		var rest transclient.Response
		var rest1 transclient.Response
		rerr := tcc.CallBranch(
			&transclient.AdjustInfo{
				UserID: req.UserId,
				Amount: req.Amount,
			},
			transRpcBusiServer+"/transclient.Trans/TransOutTry",
			transRpcBusiServer+"/transclient.Trans/TransOutConfirm",
			transRpcBusiServer+"/transclient.Trans/TransOutCancel", &rest) // 如需求拿到result， 传 reply

		rerr1 := tcc.CallBranch(
			&transclient.AdjustInfo{
				UserID: req.ToUserId,
				Amount: req.Amount,
			},
			transRpcBusiServer+"/transclient.Trans/TransInTry",
			transRpcBusiServer+"/transclient.Trans/TransInConfirm",
			transRpcBusiServer+"/transclient.Trans/TransInCancel", &rest1)
		if rerr != nil || rerr1 != nil {
			return fmt.Errorf("tcc error")
		}

		return nil
	})

	return &types.TransResponse{}, err
}
