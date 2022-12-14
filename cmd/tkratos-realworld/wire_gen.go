// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"tkratos-realworld/internal/biz"
	"tkratos-realworld/internal/conf"
	"tkratos-realworld/internal/data"
	"tkratos-realworld/internal/server"
	"tkratos-realworld/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	profileRepo := data.NewProfileRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, profileRepo, logger, jwt)
	articleRepo := data.NewArticleRepo(dataData, logger)
	commentRepo := data.NewCommentRepo(dataData, logger)
	socialUsecase := biz.NewSocialUsecase(articleRepo, profileRepo, commentRepo, logger)
	realworldService := service.NewRealworldService(userUsecase, socialUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, realworldService, logger)
	grpcServer := server.NewGRPCServer(confServer, realworldService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
