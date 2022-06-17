package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"tcc/trans/internal/svc"
	"tcc/trans/trans"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransOutTryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransOutTryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransOutTryLogic {
	return &TransOutTryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransOutTryLogic) TransOutTry(in *trans.AdjustInfo) (*trans.Response, error) {
	db, err := sqlx.NewMysql(l.svcCtx.Config.Mysql.DataSource).RawDB()
	if err != nil {
		return nil, err
	}
	// 获取子事务屏障
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, err
	}

	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 冻结
		result, err := l.svcCtx.UserAccountModel.TxAdjustTrading(tx, in.UserID, -in.Amount)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err == nil && affected == 0 {
			return fmt.Errorf("update error, balance not enough")
		}
		return err
	}); err != nil {
		return nil, err
	}

	return &trans.Response{}, nil
}
