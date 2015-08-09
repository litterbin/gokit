package main

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

type Service struct {
	Logger log.Logger
}

func (s *Service) Add(_ context.Context, a, b int64) int64 {
	s.Logger.Log("Add", "start")
	return a + b
}
