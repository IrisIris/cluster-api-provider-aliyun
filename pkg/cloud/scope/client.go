package scope

import (
	svcs "github.com/IrisIris/cluster-api-provider-aliyun/pkg/cloud/services"
)

type ACKClients struct {
	ECS svcs.ECSMachineInterface
}
