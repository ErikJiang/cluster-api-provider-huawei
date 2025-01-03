package network

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

func (s *Service) reconcileRouteTables() error {
	// Check if route table exists, if not create it
	request := &model.ListRouteTablesRequest{}
	response, err := s.vpcClient.ListRouteTables(request)
	if err != nil {
		return errors.Wrap(err, "failed to list route tables")
	}

	if len(*response.Routetables) == 0 {

		createRequest := &model.CreateRouteTableRequest{}
		nameRouteTable := "k8s-route-table"
		routetableBody := &model.CreateRouteTableReq{
			Name:  &nameRouteTable,
			VpcId: s.scope.VPC().Id,
		}
		createRequest.Body = &model.CreateRoutetableReqBody{
			Routetable: routetableBody,
		}
		_, err := s.vpcClient.CreateRouteTable(createRequest)
		if err != nil {
			return errors.Wrap(err, "failed to create route table")
		}
		klog.Infof("Created route table")
	} else {
		klog.Infof("Route table already exists")
	}

	return nil
}

func (s *Service) deleteRouteTables() error {
	request := &model.ListRouteTablesRequest{}
	response, err := s.vpcClient.ListRouteTables(request)
	if err != nil {
		return errors.Wrap(err, "failed to list route tables")
	}

	for _, routeTable := range *response.Routetables {
		deleteRequest := &model.DeleteRouteTableRequest{
			RoutetableId: routeTable.Id,
		}
		_, err := s.vpcClient.DeleteRouteTable(deleteRequest)
		if err != nil {
			return errors.Wrapf(err, "failed to delete route table %s", routeTable.Id)
		}
		klog.Infof("Deleted route table %s", routeTable.Id)
	}

	return nil
}
