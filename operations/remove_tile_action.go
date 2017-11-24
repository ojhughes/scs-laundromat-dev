package operations

import "github.com/pivotal-cf/om/api"

//import "github.com/pivotal-cf/om/api"
//import "github.com/pivotal-cf/om/commands"

type RemoveIncompleteTileService struct {
	stagedService stagedService
}

func NewRemoveIncompleteTileService(ss stagedService) RemoveIncompleteTileService {
	return RemoveIncompleteTileService{
		stagedService: ss,
	}
}

func (d RemoveIncompleteTileService) RemoveIncompleteInstall(productName string) (PostInstallState, error) {
	err := d.stagedService.Unstage(api.UnstageProductInput{ProductName: productName})
	if err != nil {
		return 0, err
	}
	return 0, nil
}

type RevertChangeService struct {
}

func (d RevertChangeService) RevertChanges(input GetStateInput) {
	//dashboardService := api.NewDashboardService(authClient)
	//revertChangeService := commands.NewRevertStagedChanges(dashboardService, logger.Info)
	//revertChangeService.Execute([]string{""})
}