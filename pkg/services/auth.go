package services

import (
	"context"
	"net/http"

	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/db"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/models"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/pb"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/utils"
)

type Server struct {
	H   db.Handler
	Jwt utils.JwtWrapper
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := s.H.DB.Limit(1).Find(&models.User{Email: req.Email}); result.RowsAffected > 0 {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "Email already exists",
		}, nil
	}

	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Password = utils.HashPassword(req.Password)

	s.H.DB.Create(&user)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Invalid login",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Invalid login",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	if result := s.H.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "Invalid login",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil

}
