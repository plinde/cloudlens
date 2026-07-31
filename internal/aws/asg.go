package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	asgTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/rs/zerolog/log"
)

func GetASGs(cfg aws.Config) ([]ASGResp, error) {
	var asgs []ASGResp
	asgClient := autoscaling.NewFromConfig(cfg)

	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(asgClient, &autoscaling.DescribeAutoScalingGroupsInput{})
	for paginator.HasMorePages() {
		result, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Info().Msg(fmt.Sprintf("Error fetching ASGs: %v", err))
			return nil, err
		}
		for _, g := range result.AutoScalingGroups {
			asg := buildASGResp(g)
			asgs = append(asgs, asg)
		}
	}
	return asgs, nil
}

func buildASGResp(g asgTypes.AutoScalingGroup) ASGResp {
	azs := azList(g.AvailabilityZones)
	tps := tpList(g.TerminationPolicies)
	lt := launchTemplate(g)

	return ASGResp{
		Name:                    safeString(g.AutoScalingGroupName),
		AutoScalingGroupName:    safeString(g.AutoScalingGroupName),
		ARN:                     safeString(g.AutoScalingGroupARN),
		AvailabilityZones:       azs,
		MinSize:                 strconv.Itoa(int(safeInt32(g.MinSize))),
		MaxSize:                 strconv.Itoa(int(safeInt32(g.MaxSize))),
		DesiredCapacity:         strconv.Itoa(int(safeInt32(g.DesiredCapacity))),
		DefaultCooldown:         strconv.Itoa(int(safeInt32(g.DefaultCooldown))),
		HealthCheckType:         safeString(g.HealthCheckType),
		HealthCheckGracePeriod:  strconv.Itoa(int(safeInt32(g.HealthCheckGracePeriod))),
		CreatedTime:             fmt.Sprintf("%v", g.CreatedTime),
		Status:                  safeString(g.Status),
		LaunchConfigurationName: safeString(g.LaunchConfigurationName),
		LaunchTemplate:          lt,
		TerminationPolicies:     tps,
	}
}

func GetSingleASG(cfg aws.Config, asgName string) string {
	asgClient := autoscaling.NewFromConfig(cfg)
	result, err := asgClient.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching ASG %s: %v", asgName, err))
		return ""
	}
	if len(result.AutoScalingGroups) == 0 {
		return ""
	}
	r, _ := json.MarshalIndent(result.AutoScalingGroups[0], "", " ")
	return string(r)
}

func UpdateASGDesiredCapacity(cfg aws.Config, asgName string, desired int32) error {
	asgClient := autoscaling.NewFromConfig(cfg)
	_, err := asgClient.SetDesiredCapacity(context.Background(), &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int32(desired),
		HonorCooldown:        aws.Bool(false),
	})
	return err
}

func UpdateASGSize(cfg aws.Config, asgName string, minSize, maxSize, desiredCapacity int32) error {
	asgClient := autoscaling.NewFromConfig(cfg)
	_, err := asgClient.UpdateAutoScalingGroup(context.Background(), &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(asgName),
		MinSize:              aws.Int32(minSize),
		MaxSize:              aws.Int32(maxSize),
		DesiredCapacity:      aws.Int32(desiredCapacity),
	})
	return err
}

func azList(zones []string) string {
	if len(zones) == 0 {
		return ""
	}
	return strings.Join(zones, ", ")
}

func tpList(policies []string) string {
	if len(policies) == 0 {
		return ""
	}
	return strings.Join(policies, ", ")
}

func launchTemplate(g asgTypes.AutoScalingGroup) string {
	if g.LaunchTemplate != nil {
		if g.LaunchTemplate.LaunchTemplateName != nil {
			return *g.LaunchTemplate.LaunchTemplateName
		}
	}
	return safeString(g.LaunchConfigurationName)
}

func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func GetASGFormatted(cfg aws.Config, asgName string) string {
	asgClient := autoscaling.NewFromConfig(cfg)
	result, err := asgClient.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil || len(result.AutoScalingGroups) == 0 {
		return ""
	}
	g := result.AutoScalingGroups[0]
	return fmt.Sprintf(`Name:                     %s
ARN:                      %s
Status:                   %s
Min Size:                 %d
Max Size:                 %d
Desired Capacity:         %d
Default Cooldown:         %d
Health Check Type:        %s
Health Check Grace Period: %d
Availability Zones:       %s
Launch Configuration:     %s
Launch Template:          %s
Termination Policies:     %s
Created:                  %v
Instances:                %d`,
		safeString(g.AutoScalingGroupName),
		safeString(g.AutoScalingGroupARN),
		safeString(g.Status),
		safeInt32(g.MinSize),
		safeInt32(g.MaxSize),
		safeInt32(g.DesiredCapacity),
		safeInt32(g.DefaultCooldown),
		safeString(g.HealthCheckType),
		safeInt32(g.HealthCheckGracePeriod),
		azList(g.AvailabilityZones),
		safeString(g.LaunchConfigurationName),
		launchTemplate(g),
		tpList(g.TerminationPolicies),
		fmt.Sprintf("%v", g.CreatedTime),
		len(g.Instances),
	)
}

func GetASGInstances(cfg aws.Config, asgName string) ([]ASGInstanceResp, error) {
	asgClient := autoscaling.NewFromConfig(cfg)
	result, err := asgClient.DescribeAutoScalingGroups(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{asgName},
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error fetching ASG instances for %s: %v", asgName, err))
		return nil, err
	}
	if len(result.AutoScalingGroups) == 0 {
		return nil, nil
	}

	instances := result.AutoScalingGroups[0].Instances
	resp := make([]ASGInstanceResp, 0, len(instances))
	for _, inst := range instances {
		protected := "No"
		if inst.ProtectedFromScaleIn != nil && *inst.ProtectedFromScaleIn {
			protected = "Yes"
		}
		resp = append(resp, ASGInstanceResp{
			InstanceId:       safeString(inst.InstanceId),
			InstanceType:     safeString(inst.InstanceType),
			AvailabilityZone: safeString(inst.AvailabilityZone),
			LifecycleState:   string(inst.LifecycleState),
			HealthStatus:     safeString(inst.HealthStatus),
			Protected:        protected,
		})
	}
	return resp, nil
}
