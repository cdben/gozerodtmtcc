package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tcc/trans/internal/svc"
	"tcc/trans/trans"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransInConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransInConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransInConfirmLogic {
	return &TransInConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransInConfirmLogic) TransInConfirm(in *trans.AdjustInfo) (*trans.Response, error) {
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	// 获取子事务屏障
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 冻结
		result, err := l.svcCtx.UserAccountModel.TxAdjustBalance(tx, in.UserID, in.Amount)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err == nil && affected == 0 {
			return dtmcli.ErrFailure
		}
		return err
	}); err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &trans.Response{}, nil
}
