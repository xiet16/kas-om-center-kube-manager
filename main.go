package main

import (
	"fmt"
	"kas/om-center/kube-manager/client_example"
)

func main() {
	//client_example.DeletePod()
	//client_example.GetPod("")
	dm := client_example.NewDeploymentManager()
	serverName := "bc-order"
	dm.DeleteDeploy(fmt.Sprintf("deploy-%s", serverName))
	dm.CreateDemoDeploy(serverName, "latest")
	sm := client_example.NewServiceManager("")
	sm.CreateService(serverName)
}
