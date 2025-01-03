package securitygroup

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"k8s.io/klog/v2"
)

func (s *Service) ReconcileSecurityGroups() error {
	klog.Info("Reconciling security groups")

	// Create a basic security group
	securityGroupName := "k8s-basic-sg"
	createSecurityGroupRequest := &model.CreateSecurityGroupRequest{
		Body: &model.CreateSecurityGroupRequestBody{
			SecurityGroup: &model.CreateSecurityGroupOption{
				Name: securityGroupName,
			},
		},
	}

	createSecurityGroupResponse, err := s.vpcClient.CreateSecurityGroup(createSecurityGroupRequest)
	if err != nil {
		return fmt.Errorf("failed to create security group: %v", err)
	}

	securityGroupID := createSecurityGroupResponse.SecurityGroup.Id
	klog.Infof("Created security group: %s", securityGroupID)

	ingressRules := []model.NeutronCreateSecurityGroupRuleOption{}
	protocol := "tcp"
	remoteIpPrefix := "0.0.0.0/0"
	portRangeMin, portRangeMax := int32(22), int32(22)
	ingressRules = append(ingressRules, model.NeutronCreateSecurityGroupRuleOption{
		SecurityGroupId: securityGroupID,
		Direction:       model.GetNeutronCreateSecurityGroupRuleOptionDirectionEnum().INGRESS,
		PortRangeMax:    &portRangeMax,
		PortRangeMin:    &portRangeMin,
		Protocol:        &protocol,
		RemoteIpPrefix:  &remoteIpPrefix,
	})

	portRangeMin, portRangeMax = int32(6443), int32(6443)
	ingressRules = append(ingressRules, model.NeutronCreateSecurityGroupRuleOption{
		SecurityGroupId: securityGroupID,
		Direction:       model.GetNeutronCreateSecurityGroupRuleOptionDirectionEnum().INGRESS,
		PortRangeMax:    &portRangeMax,
		PortRangeMin:    &portRangeMin,
		Protocol:        &protocol,
		RemoteIpPrefix:  &remoteIpPrefix,
	})

	for _, rule := range ingressRules {
		createSecurityGroupRuleRequest := &model.NeutronCreateSecurityGroupRuleRequest{
			Body: &model.NeutronCreateSecurityGroupRuleRequestBody{
				SecurityGroupRule: &rule,
			},
		}

		_, err := s.vpcClient.NeutronCreateSecurityGroupRule(createSecurityGroupRuleRequest)
		if err != nil {
			return fmt.Errorf("failed to create security group rule: %v", err)
		}
		klog.Infof("Created security group rule: %+v", rule)
	}

	return nil
}

func (s *Service) DeleteSecurityGroups() error {
	klog.Info("Deleting security groups")

	// Retrieve the security group by name
	securityGroupName := "k8s-basic-sg"
	listSecurityGroupsRequest := &model.NeutronListSecurityGroupsRequest{
		Name: &securityGroupName,
	}

	listSecurityGroupsResponse, err := s.vpcClient.NeutronListSecurityGroups(listSecurityGroupsRequest)
	if err != nil {
		return fmt.Errorf("failed to list security groups: %v", err)
	}

	securityGroups := *listSecurityGroupsResponse.SecurityGroups
	if len(securityGroups) == 0 {
		klog.Infof("No security group found with name: %s", securityGroupName)
		return nil
	}

	securityGroupID := securityGroups[0].Id
	klog.Infof("Found security group: %s", securityGroupID)

	// Delete all security group rules
	listSecurityGroupRulesRequest := &model.NeutronListSecurityGroupRulesRequest{
		SecurityGroupId: &securityGroupID,
	}

	listSecurityGroupRulesResponse, err := s.vpcClient.NeutronListSecurityGroupRules(listSecurityGroupRulesRequest)
	if err != nil {
		return fmt.Errorf("failed to list security group rules: %v", err)
	}

	for _, rule := range *listSecurityGroupRulesResponse.SecurityGroupRules {
		deleteSecurityGroupRuleRequest := &model.NeutronDeleteSecurityGroupRuleRequest{
			SecurityGroupRuleId: rule.Id,
		}

		_, err := s.vpcClient.NeutronDeleteSecurityGroupRule(deleteSecurityGroupRuleRequest)
		if err != nil {
			return fmt.Errorf("failed to delete security group rule: %v", err)
		}
		klog.Infof("Deleted security group rule: %s", rule.Id)
	}

	// Delete the security group
	deleteSecurityGroupRequest := &model.NeutronDeleteSecurityGroupRequest{
		SecurityGroupId: securityGroupID,
	}

	_, err = s.vpcClient.NeutronDeleteSecurityGroup(deleteSecurityGroupRequest)
	if err != nil {
		return fmt.Errorf("failed to delete security group: %v", err)
	}

	klog.Infof("Deleted security group: %s", securityGroupID)
	return nil
}
