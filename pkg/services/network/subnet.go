package network

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

func (s *Service) reconcileSubnets() error {
	// Check if subnet exists, if not create it
	request := &model.ListSubnetsRequest{}
	response, err := s.vpcClient.ListSubnets(request)
	if err != nil {
		return errors.Wrap(err, "failed to list subnets")
	}

	if len(*response.Subnets) == 0 {
		createRequest := &model.CreateSubnetRequest{}
		subnetbody := &model.CreateSubnetOption{
			Name:      "k8s-subnet",
			Cidr:      "192.168.1.0/24",
			VpcId:     s.scope.VPC().Id,
			GatewayIp: "192.168.1.1",
		}
		createRequest.Body = &model.CreateSubnetRequestBody{
			Subnet: subnetbody,
		}
		_, err := s.vpcClient.CreateSubnet(createRequest)
		if err != nil {
			return errors.Wrap(err, "failed to create subnet")
		}
		klog.Infof("Created subnet")
	} else {
		klog.Infof("Subnet already exists")
	}

	return nil
}

func (s *Service) deleteSubnets() error {
	request := &model.ListSubnetsRequest{}
	response, err := s.vpcClient.ListSubnets(request)
	if err != nil {
		return errors.Wrap(err, "failed to list subnets")
	}

	for _, subnet := range *response.Subnets {
		deleteRequest := &model.DeleteSubnetRequest{
			VpcId:    subnet.VpcId,
			SubnetId: subnet.Id,
		}
		_, err := s.vpcClient.DeleteSubnet(deleteRequest)
		if err != nil {
			return errors.Wrapf(err, "failed to delete subnet %s", subnet.Id)
		}
		klog.Infof("Deleted subnet %s", subnet.Id)
	}

	return nil
}
