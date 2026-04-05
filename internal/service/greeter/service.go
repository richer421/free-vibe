package greeter

import (
	"context"

	v1 "free-vibe-coding/api/helloworld/v1"
	greeterbiz "free-vibe-coding/internal/biz/greeter"
)

// Service is a greeter service.
type Service struct {
	v1.UnimplementedGreeterServer

	uc *greeterbiz.Usecase
}

// NewService new a greeter service.
func NewService(uc *greeterbiz.Usecase) *Service {
	return &Service{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *Service) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.Create(ctx, &greeterbiz.Entity{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
