package service

import "KillerFeature/ClientSide/internal"

type Service struct {
}

func NewService() internal.Usecase {
	return &Service{}
}
