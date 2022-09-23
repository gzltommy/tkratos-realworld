package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "tkratos-realworld/api/realworld/v1"
	"tkratos-realworld/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewRealworldService)

type RealworldService struct {
	v1.UnimplementedRealworldServer

	uc  *biz.UserUsecase
	sc  *biz.SocialUsecase
	log *log.Helper
}

func NewRealworldService(uc *biz.UserUsecase, sc *biz.SocialUsecase, logger log.Logger) *RealworldService {
	return &RealworldService{uc: uc, sc: sc, log: log.NewHelper(logger)}
}
