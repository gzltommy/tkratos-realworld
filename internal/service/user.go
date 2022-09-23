package service

import (
	"context"
	"tkratos-realworld/internal/biz"

	v1 "tkratos-realworld/api/realworld/v1"
)

func (s *RealworldService) Login(ctx context.Context, req *v1.LoginRequest) (reply *v1.UserReply, err error) {
	rv, err := s.uc.Login(ctx, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Username: rv.Username,
			Token:    rv.Token,
		},
	}, nil
}

func (s *RealworldService) Register(ctx context.Context, req *v1.RegisterRequest) (reply *v1.UserReply, err error) {
	u, err := s.uc.Register(ctx, req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Email:    u.Email,
			Username: u.Username,
			Token:    u.Token,
		},
	}, nil
}

func (s *RealworldService) GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (reply *v1.UserReply, err error) {
	u, err := s.uc.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Username: u.Username,
			Image:    u.Image,
			Bio:      u.Bio,
		},
	}, nil
}

func (s *RealworldService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (reply *v1.UserReply, err error) {
	u, err := s.uc.UpdateUser(ctx, &biz.UserUpdate{
		Email:    req.User.GetEmail(),
		Username: req.User.GetUsername(),
		Password: req.User.GetPassword(),
		Bio:      req.User.GetBio(),
		Image:    req.User.GetImage(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserReply{
		User: &v1.UserReply_User{
			Username: u.Username,
			Email:    u.Email,
			Image:    u.Image,
			Bio:      u.Bio,
		},
	}, nil
}
