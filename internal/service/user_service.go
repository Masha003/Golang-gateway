package service

import (
	"context"

	"github.com/Masha003/Golang-gateway/internal/config"
	"github.com/Masha003/Golang-gateway/internal/models"
	"github.com/Masha003/Golang-gateway/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserService interface {
	Register(user models.RegisterUser) (models.Token, error)
	Login(user models.LoginUser) (models.Token, error)
	FindAll(query models.PaginationQuery) ([]models.User, error)
	FindById(id string) (models.User, error)
	Delete(id string) error
	Close() error
}

func NewUserService(cfg config.Config) (UserService, error) {
	conn, err := grpc.Dial(cfg.GrpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)

	us := &userService{
		cfg:    cfg,
		conn:   conn,
		client: client,
	}

	return us, nil
}

type userService struct {
	cfg    config.Config
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func mapUser(pbUser *pb.User) models.User {
	var user models.User

	user.ID = pbUser.Id
	user.Email = pbUser.Email
	user.Name = pbUser.Name
	user.Image = pbUser.Image
	user.CreatedAt = pbUser.CreatedAt.AsTime()
	user.UpdatedAt = pbUser.UpdatedAt.AsTime()
	user.DeletedAt = pbUser.DeletedAt.AsTime()

	return user
}

func (s *userService) Close() error {
	return s.conn.Close()
}

func (s *userService) Login(loginUser models.LoginUser) (models.Token, error) {
	resp, err := s.client.Login(context.Background(), &pb.LoginUser{
		Email:    loginUser.Email,
		Password: loginUser.Password,
	})

	var token models.Token

	if err != nil {
		return token, err
	}

	token.RefreshToken = resp.RefreshToken
	token.Token = resp.Token
	token.User = mapUser(resp.User)

	return token, nil
}

func (s *userService) FindAll(query models.PaginationQuery) ([]models.User, error) {
	resp, err := s.client.GetAll(context.Background(), &pb.UsersQuery{
		Page: int64(query.Page),
		Size: int64(query.Size),
	})

	var users []models.User

	if err != nil {
		return users, err
	}

	for _, pbuser := range resp.Users {
		users = append(users, mapUser(pbuser))
	}

	return users, nil
}

func (s *userService) FindById(id string) (models.User, error) {
	resp, err := s.client.GetById(context.Background(), &pb.UserId{
		Id: id,
	})

	var user models.User

	if err != nil {
		return user, nil
	}

	user = mapUser(resp)

	return user, nil
}

func (s *userService) Register(registerUser models.RegisterUser) (models.Token, error) {
	resp, err := s.client.Register(context.Background(), &pb.RegisterUser{
		Email:    registerUser.Email,
		Name:     registerUser.Name,
		Password: registerUser.Password,
	})

	var token models.Token

	if err != nil {
		return token, err
	}

	token.RefreshToken = resp.RefreshToken
	token.Token = resp.Token
	token.User = mapUser(resp.User)

	return token, nil
}

func (s *userService) Delete(id string) error {
	_, err := s.client.Delete(context.Background(), &pb.UserId{
		Id: id,
	})
	return err
}
