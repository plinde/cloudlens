package aws

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetInstancesBuildsCompatibleSDKMiddleware(t *testing.T) {
	called := false
	cfg := aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("test", "test", ""),
		HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			called = true
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(
					`<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2016-11-15/"><requestId>test</requestId><reservationSet/></DescribeInstancesResponse>`,
				)),
				Request: req,
			}, nil
		}),
	}

	instances, err := GetInstances(cfg)
	if err != nil {
		t.Fatalf("GetInstances returned middleware error: %v", err)
	}
	if !called {
		t.Fatal("expected EC2 request to reach HTTP transport")
	}
	if len(instances) != 0 {
		t.Fatalf("expected no instances, got %d", len(instances))
	}
}

type Ec2API interface {
	GetEc2Instances(ctx context.Context, params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)
}

func GetEc2InstancesTest(ctx context.Context, api Ec2API) (*ec2.DescribeInstancesOutput, error) {
	object, err := api.GetEc2Instances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

type mockGetEc2Instances func(ctx context.Context, params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error)

func (m mockGetEc2Instances) GetEc2Instances(ctx context.Context, params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return m(ctx, params)
}

func TestGetec2Instances(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) Ec2API
		expect ec2.DescribeInstancesOutput
	}{
		{
			client: func(t *testing.T) Ec2API {
				return mockGetEc2Instances(func(ctx context.Context, params *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
					t.Helper()
					return &ec2.DescribeInstancesOutput{Reservations: []types.Reservation{{Instances: []types.Instance{{InstanceId: aws.String("ec2-instance-1"), InstanceType: types.InstanceType(*aws.String("t2.micro"))}}}}}, nil

				})
			},
			expect: ec2.DescribeInstancesOutput{Reservations: []types.Reservation{{Instances: []types.Instance{{InstanceId: aws.String("ec2-instance-1"), InstanceType: types.InstanceType(*aws.String("t2.micro"))}}}}},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			got, err := GetEc2InstancesTest(ctx, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			for i := 0; i < len(got.Reservations); i++ {
				for j := 0; j < len(got.Reservations[i].Instances); j++ {
					fmt.Println("got:", *got.Reservations[i].Instances[j].InstanceId)
					fmt.Println("expect:", *tt.expect.Reservations[i].Instances[j].InstanceId)
					if *got.Reservations[i].Instances[j].InstanceId != *tt.expect.Reservations[i].Instances[j].InstanceId {
						t.Errorf("expect %v, got %v", *tt.expect.Reservations[i].Instances[j].InstanceId, *got.Reservations[i].Instances[j].InstanceId)
					}
				}
			}
		})
	}
}
