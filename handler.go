package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/eluts15/shipping-user-service/proto/auth"
	"golang.org/x/crypto/bcrypt"
)

const topic = "user.created"

type service struct {
	repo         Repository
	tokenService Authable
}

func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (srv *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with: ", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	log.Println(user, err)
	if err != nil {
		return err
	}

	// Compares the given password to the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (srv *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	log.Println("Creating user: ", req)

	// Generate hashed password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("error hashing password: %v", err))
	}

	token, err := srv.tokenService.Encode(req)
	if err != nil {
		return err
	}

	res.User = req
	res.Token = &pb.Token{Token: token}

	return nil
}

func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	// Decode token
	claims, err := srv.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}
