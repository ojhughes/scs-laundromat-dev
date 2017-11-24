package main

import (
	"time"

	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"github.com/pivotal-cf/scs-laundromat-dev/operations"
	"sync"
	"github.com/pivotal-cf/om/commands"
)



func main() {

	//target := os.Getenv("TARGET")
	//opsManagerUser := os.Getenv("OPS_MAN_USER")
	//opsManagerPassword := os.Getenv("OPS_MAN_PASSWORD")
	tileName := "p-spring-cloud-services"
	target := "pcf.indigo.springapps.io"
	opsManagerUser := "pivotalcf"
	opsManagerPassword := "pivotalcf"

	requestTimeout := time.Duration(1800) * time.Second
	authClient, err := network.NewOAuthClient(target, opsManagerUser, opsManagerPassword, "", "", true, false, requestTimeout)
	if err != nil {
		logger.Error.Fatal(err)
	}
	var wg sync.WaitGroup
	errandDisabler := api.NewErrandsService(authClient)
	disableErrandService := operations.NewDisableErrandService(errandDisabler)
	go disableErrandService.DisableSystemErrands(&wg)

	stageService := api.NewStagedProductsService(authClient)
	//unstageInput := api.UnstageProductInput{ProductName: tileName}
	pendingChangesService := api.NewPendingChangesService(authClient)
	deployedProductService := api.NewDeployedProductsService(authClient)

	determineState := operations.NewDetermineTileStateService(stageService, deployedProductService, pendingChangesService)
	tileState, err := determineState.PreInstallState(operations.GetStateInput{ProductName: tileName})
	waitTimeout(&wg, 20 * time.Second)
	switch tileState {
	case operations.StagedButNotDeployed: {
		removeIncompleteTileService := operations.NewRemoveIncompleteTileService(stageService)
		removeIncompleteTileService.RemoveIncompleteInstall(tileName)
	}
	case operations.PendingUpdate: {
		dashboardService := api.NewDashboardService(authClient)
		revertChangeService := operations.NewRevertChangeService(dashboardService)
		revertChangeService.RevertChanges()
	}
	case operations.PendingDeletion: {

	}
	case operations.StagedAndDeployed: {

	}
	}
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}