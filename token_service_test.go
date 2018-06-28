package main

import (
	"strings"
	"testing"

	pb "github.com/eluts15/shipping-user-service/proto/auth"
)

var (
	user = &pb.User{
		Id:    "123",
		Email: "eluts15dev@gmail.com",
	}
)

type MockRepo struct{}

func (repo *MockRepo) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	return users, nil
}

func (repo *MockRepo) Get(id string) (*pb.User, error) {
	var user *pb.User
	return user, nil
}

func (repo *MockRepo) Create(user *pb.User) error {
	return nil
}

func (repo *MockRepo) GetByEmail(email string) (*pb.User, error) {
	var user *pb.User
	return user, nil
}

func newInstance() Authable {
	repo := &MockRepo{}
	return &TokenService{repo}
}

func TestCanCreateToken(t *testing.T) {
	srv := newInstance()
	token, err := srv.Encode(user)
	if err != nil {
		t.Fail()
	}

	if token == "" {
		t.Fail()
	}

	if len(string.Split(token, ".")) != 3 {
		t.Fail()
	}
}

func TestCanDecodeToken(t *testing.T) {
	srv := newInstance()
	token, err := srv.Encode(user)
	if err != nil {
		t.Fail()
	}
	claims, err := srv.Decode(token)
	if err != nil {
		t.Fail()
	}
	if claims.User == nil {
		t.Fail()
	}
	if claims.User.Email != "eluts15dev@gmail.com" {
		t.Fail()
	}

}

func TestThrowsErrorIfTokenInvalid(t *testing.T) {
	srv := newInstance()
	_, err := srv.Decode("xxx.xxx.xxx")
	if err == nil {
		t.Fail()
	}
}
