package main

import (
	"time"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
	"github.com/pivotal-cf/scs-laundromat-dev/commands"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"os"
)

const productName = "p-spring-cloud-services"

func main() {

	target := os.Getenv("TARGET")
	opsManagerUser := os.Getenv("OPS_MAN_USER")
	opsManagerPassword := os.Getenv("OPS_MAN_PASSWORD")

	allowedArgs := make(map[string]bool)
	allowedArgs["remove"] = true
	allowedArgs["clean-install"] = true
	if len(os.Args) != 2 {
		logger.Error.Fatalf("Usage: %s [remove|clean-install]", os.Args[0])
	}
	command := os.Args[1]
	if _, ok := allowedArgs[command]; !ok {
		logger.Error.Fatalf("Usage: %s [remove|clean-install]", os.Args[0])
	}

	requestTimeout := time.Duration(1800) * time.Second
	authClient, err := network.NewOAuthClient(target, opsManagerUser, opsManagerPassword, "", "", true, false, requestTimeout)
	if err != nil {
		logger.Error.Fatal(err)
	}
	//Wire up dependencies
	errandDisabler := api.NewErrandsService(authClient)
	stageService := api.NewStagedProductsService(authClient)
	pendingChangesService := api.NewPendingChangesService(authClient)
	deployedProductService := api.NewDeployedProductsService(authClient)
	installationsService := api.NewInstallationsService(authClient)
	credentialService := api.NewCredentialsService(authClient, progress.NewBar())

	if command == "remove" {
		removeTileService := commands.NewRemoveTileService(
			errandDisabler,
			stageService,
			pendingChangesService,
			deployedProductService,
			installationsService,
			credentialService,
		)
		removeTileService.RemoveTile(productName, target)
	} else if command == "clean-install" {
		logger.Error.Fatal("clean-install not yet implemented")

	}
}
