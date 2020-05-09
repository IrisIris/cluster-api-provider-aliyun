package ecs

import (
	"github.com/IrisIris/cluster-api-provider-aliyun/api/v1alpha3"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// use sdk to create ecs instance
func (s *Service) RunInstances(request *ecs.RunInstancesRequest) (response *ecs.RunInstancesResponse, err error) {
	return response, nil
}

func (s *Service) CreateInstances() (*v1alpha3.Instance, error) {
	input := v1alpha3.Instance{}
	// make run instance request
	// use SDK to run Instance
	createRequest := ecs.CreateRunInstancesRequest()
	// get client

	return nil, nil
}
