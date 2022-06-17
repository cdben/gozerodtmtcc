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

type TransOutConfirmLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransOutConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransOutConfirmLogic {
	return &TransOutConfirmLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransOutConfirmLogic) TransOutConfirm(in *trans.AdjustInfo) (*trans.Response, error) {
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
		result, err := l.svcCtx.UserAccountModel.TxAdjustBalance(tx, in.UserID, -in.Amount)
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
