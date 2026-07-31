package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ASG struct {
	Accessor
	ctx context.Context
}

func (a *ASG) Init(ctx context.Context) {
	a.ctx = ctx
}

func (a *ASG) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	asgs, err := aws.GetASGs(cfg)
	objs := make([]Object, len(asgs))
	for i, obj := range asgs {
		objs[i] = obj
	}
	return objs, err
}

func (a *ASG) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (a *ASG) Describe(asgName string) (string, error) {
	cfg, ok := a.ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	res := aws.GetASGFormatted(cfg, asgName)
	return res, nil
}

func (a *ASG) UpdateSize(ctx context.Context, asgName string, minSize, maxSize, desiredCapacity int32) error {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	return aws.UpdateASGSize(cfg, asgName, minSize, maxSize, desiredCapacity)
}

func (a *ASG) UpdateDesiredCapacity(ctx context.Context, asgName string, desired int32) error {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	return aws.UpdateASGDesiredCapacity(cfg, asgName, desired)
}
