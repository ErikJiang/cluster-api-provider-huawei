package network

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

func (s *Service) reconcilePublicIPs() error {
	// Check if internet gateway exists, if not create it
	request := &model.ListPublicipsRequest{}
	response, err := s.eipClient.ListPublicips(request)
	if err != nil {
		return errors.Wrap(err, "failed to list public ips")
	}

	if len(*response.Publicips) == 0 {
		createRequest := &model.CreatePublicipRequest{}
		publicipbody := &model.CreatePublicipOption{
			Type: "5_bgp",
		}
		createRequest.Body = &model.CreatePublicipRequestBody{
			Publicip: publicipbody,
		}
		_, err := s.eipClient.CreatePublicip(createRequest)
		if err != nil {
			return errors.Wrap(err, "failed to create public ip")
		}
		klog.Infof("Created public ip")
	} else {
		klog.Infof("Public ip already exists")
	}

	return nil
}

func (s *Service) deletePublicIPs() error {
	request := &model.ListPublicipsRequest{}
	response, err := s.eipClient.ListPublicips(request)
	if err != nil {
		return errors.Wrap(err, "failed to list public ips")
	}

	for _, publicip := range *response.Publicips {
		deleteRequest := &model.DeletePublicipRequest{
			PublicipId: *publicip.Id,
		}
		_, err := s.eipClient.DeletePublicip(deleteRequest)
		if err != nil {
			return errors.Wrapf(err, "failed to delete public ip %s", *publicip.Id)
		}
		klog.Infof("Deleted public ip %s", *publicip.Id)
	}

	return nil
}
