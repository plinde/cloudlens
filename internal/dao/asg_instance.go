package dao

import (
	"context"
	"fmt"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/aws"
	"github.com/rs/zerolog/log"
)

type ASGInstance struct {
	Accessor
	ctx context.Context
}

func (a *ASGInstance) Init(ctx context.Context) {
	a.ctx = ctx
}

func (a *ASGInstance) List(ctx context.Context) ([]Object, error) {
	cfg, ok := ctx.Value(internal.KeySession).(awsV2.Config)
	if !ok {
		log.Err(fmt.Errorf("conversion err: Expected awsV2.Config but got %v", cfg))
	}
	asgName := fmt.Sprintf("%v", ctx.Value(internal.AsgName))
	insts, err := aws.GetASGInstances(cfg, asgName)
	objs := make([]Object, len(insts))
	for i, obj := range insts {
		objs[i] = obj
	}
	return objs, err
}

func (a *ASGInstance) Get(ctx context.Context, path string) (Object, error) {
	return nil, nil
}

func (a *ASGInstance) Describe(instanceId string) (string, error) {
	return "", nil
}
