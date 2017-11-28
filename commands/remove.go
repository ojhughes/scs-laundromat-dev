package commands

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"github.com/pivotal-cf/scs-laundromat-dev/operations"
	"sync"
	"time"
)

func NewRemoveTileService(ed api.ErrandsService,
	ss api.StagedProductsService,
	pcs api.PendingChangesService,
	dps api.DeployedProductsService,
	is api.InstallationsService,
	cs api.CredentialsService) RemoveTileService {
	return RemoveTileService{
		ed,
		ss,
		pcs,
		dps,
		is,
		cs,
	}
}

type RemoveTileService struct {
	errandDisabler         api.ErrandsService
	stageService           api.StagedProductsService
	pendingChangesService  api.PendingChangesService
	deployedProductService api.DeployedProductsService
	installationsService   api.InstallationsService
	credentialService      api.CredentialsService
}

func (rts RemoveTileService) RemoveTile(productName string, target string) {
	//setup internal interfaces
	performInstallService := operations.NewPerformInstallService(rts.installationsService)
	removeIncompleteTileService := operations.NewRemoveIncompleteTileService(rts.stageService)
	extractCreds := operations.NewExtractCredentialService(rts.credentialService, rts.deployedProductService)
	determineState := operations.NewDetermineTileStateService(rts.stageService, rts.deployedProductService, rts.pendingChangesService)

	//Disable Errands in the background
	var wg sync.WaitGroup
	disableErrandService := operations.NewDisableErrandService(rts.errandDisabler, rts.stageService)
	go disableErrandService.DisableSystemErrands(&wg)
	wg.Add(1)

	////Get the current state of the installed tile
	tileState, err := determineState.PreInstallState(operations.GetStateInput{ProductName: productName})
	if err != nil {
		logger.Error.Fatalln(err)
	}

	//Wait for Errand states to be set before applying changes
	if waitTimeout(&wg, 2*time.Minute) {
		logger.Error.Fatalln("exceeded timeout waiting for errand state to be set")
	}
	var installOutput api.InstallationsServiceOutput
	switch tileState {
	case operations.StagedButNotDeployed:
		{
			removeIncompleteTileService.RemoveInstall(productName)
		}
	case operations.PendingUpdate:
		{
			removeIncompleteTileService.RemoveInstall(productName)
			installOutput, err = performInstallService.PerformInstall(productName)
		}
	case operations.PendingDeletion:
		{
			installOutput, err = performInstallService.PerformInstall(productName)

		}
	case operations.StagedAndDeployed:
		{
			removeIncompleteTileService.RemoveInstall(productName)
			installOutput, err = performInstallService.PerformInstall(productName)
		}
	}
	if installOutput.Status == api.StatusFailed {
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
		err = forcedUninstaller.ForceUninstall(productName)
		if err != nil {
			logger.Error.Fatalln(err)
		}
		performInstallService.PerformInstall(productName)
		if err != nil {
			logger.Error.Fatalln(err)
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
