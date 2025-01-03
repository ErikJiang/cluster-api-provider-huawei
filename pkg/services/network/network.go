package network

import "k8s.io/klog/v2"

func (s *Service) ReconcileNetwork() error {
	klog.Infof("Reconciling network")

	// VPC
	if err := s.reconcileVPC(); err != nil {
		klog.Errorf("Failed to reconcile VPC: %v", err)
		return err
	}

	// Subnets
	if err := s.reconcileSubnets(); err != nil {
		klog.Errorf("Failed to reconcile subnets: %v", err)
		return err
	}

	// Public IPs
	if err := s.reconcilePublicIPs(); err != nil {
		klog.Errorf("Failed to reconcile internet gateways: %v", err)
		return err
	}

	// Routing tables
	if err := s.reconcileRouteTables(); err != nil {
		klog.Errorf("Failed to reconcile route tables: %v", err)
		return err
	}

	klog.Infof("Reconcile network completed successfully")
	return nil
}

func (s *Service) DeleteNetwork() error {
	klog.Infof("Deleting network")

	// Delete Subnets
	if err := s.deleteSubnets(); err != nil {
		klog.Errorf("Failed to delete subnets: %v", err)
		return err
	}

	// Delete Public IPs
	if err := s.deletePublicIPs(); err != nil {
		klog.Errorf("Failed to delete internet gateways: %v", err)
		return err
	}

	// Delete Route Tables
	if err := s.deleteRouteTables(); err != nil {
		klog.Errorf("Failed to delete route tables: %v", err)
		return err
	}

	// Delete VPC
	if err := s.deleteVPC(); err != nil {
		klog.Errorf("Failed to delete VPC: %v", err)
		return err
	}

	klog.Infof("Delete network completed successfully")
	return nil
}
