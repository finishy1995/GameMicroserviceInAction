package logic

import (
	"ProjectX/base"
	"ProjectX/service/account/consts"
	"ProjectX/service/account/internal/svc"
	"ProjectX/service/account/pb/account"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrCreateAccountLogic struct {
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	logx.Logger // TODO: 需要替换为 library
}

func NewGetOrCreateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrCreateAccountLogic {
	return &GetOrCreateAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrCreateAccount 获取或新账号创建
func (l *GetOrCreateAccountLogic) GetOrCreateAccount(in *account.GetOrCreateAccountRequest) (*account.GetOrCreateAccountResponse, error) {
	// 1. 使用账号id作为索引查询数据库，如果查不到走2，如果能查到走3，如果数据库响应异常走4
	info := &account.AccountInfo{
		Nickname: in.Id,
	}
	// 2. 数据库查不到，说明这是一个新账号，创建新账号，并返回
	return &account.GetOrCreateAccountResponse{
		ErrorCode: base.ErrorCodeOK,
		IsCreated: true,
		Account:   info,
	}, nil

	// 3. 数据库能查到，比对密码，根据密码是否正确返回
	if checkPassword(in.Password, "") {
	} else {
		return &account.GetOrCreateAccountResponse{
			ErrorCode: base.ErrorCodeAccountIdPasswordWrong,
			IsCreated: false,
			Account:   info,
		}, nil
	}

	// 4. 数据库异常，接口调用错误
	l.Errorf("GetOrCreateAccount failed due to database error, id: %s, error: ", in.Id)
	return &account.GetOrCreateAccountResponse{}, consts.ErrorDatabaseUnknownErr
}

func checkPassword(userInput string, databaseStore string) bool {
	// TODO: 加密并比对密码
	return true
}
