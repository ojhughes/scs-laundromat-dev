package main

import (
	"time"

	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"github.com/pivotal-cf/scs-laundromat-dev/operations"
	"sync"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/om/progress"
	"go/token"
)



func main() {

	//target := os.Getenv("TARGET")
	//opsManagerUser := os.Getenv("OPS_MAN_USER")
	//opsManagerPassword := os.Getenv("OPS_MAN_PASSWORD")
	productName := "p-spring-cloud-services"
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
	wg.Add(1)
	stageService := api.NewStagedProductsService(authClient)
	//unstageInput := api.UnstageProductInput{ProductName: productName}
	pendingChangesService := api.NewPendingChangesService(authClient)
	deployedProductService := api.NewDeployedProductsService(authClient)
	performInstallService := operations.NewPerformInstallService(api.NewInstallationsService(authClient))
	removeIncompleteTileService := operations.NewRemoveIncompleteTileService(stageService)
	credentialService := api.NewCredentialsService(authClient, progress.NewBar())
	determineState := operations.NewDetermineTileStateService(stageService, deployedProductService, pendingChangesService)
	tileState, err := determineState.PreInstallState(operations.GetStateInput{ProductName: productName})
	if waitTimeout(&wg, 20 * time.Second) {
		logger.Error.Fatalln("exceeded timeout waiting for errand state to be set")
	}
	var installOutput api.InstallationsServiceOutput
	switch tileState {
	case operations.StagedButNotDeployed: {
		removeIncompleteTileService.RemoveInstall(productName)
	}
	case operations.PendingUpdate: {
		dashboardService := api.NewDashboardService(authClient)
		revertChangeService := operations.NewRevertChangeService(dashboardService)
		err := revertChangeService.RevertChanges()
		if err != nil {
			//Still Try and unstage tile if reverting fails
			logger.Error.Println(err)
		}
		removeIncompleteTileService.RemoveInstall(productName)
		installOutput, err = performInstallService.PerformInstall(productName)
	}
	case operations.PendingDeletion: {
		installOutput, err = performInstallService.PerformInstall(productName)

	}
	case operations.StagedAndDeployed: {
		removeIncompleteTileService.RemoveInstall(productName)
		installOutput, err = performInstallService.PerformInstall(productName)
	}
	if installOutput.Status == api.StatusFailed {
		extractCreds := operations.NewExtractCredentialService(credentialService, deployedProductService)
		cfAdminCreds, err := extractCreds.ExtractCfPassword(productName)
		if err != nil {
			logger.Error.Fatalln(err)
		}
		cfClientFactory := operations.NewCfClientFactory()
		config := &cfclient.Config{
			ApiAddress:        "https://" + target,
			Username:          cfAdminCreds["identity"],
			Password:          cfAdminCreds["password"],
			SkipSslValidation: true,
		}
		cfClientService, err := cfClientFactory.NewClient(config)
		if err != nil {
			logger.Error.Fatalln(err)
		}
		forcedUninstaller := operations.NewForceUninstallService(cfClientService)
		forcedUninstaller.ForceUninstall(productName)
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