package operations

import (
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"github.com/pivotal-cf/om/commands"
	"io/ioutil"
	"github.com/pivotal-cf/om/network"
)

type RemoveIncompleteTileService struct {
	stagedService stagedService
}

func NewRemoveIncompleteTileService(ss stagedService) RemoveIncompleteTileService {
	return RemoveIncompleteTileService{
		stagedService: ss,
	}
}

func (d RemoveIncompleteTileService) RemoveIncompleteInstall(productName string) error {
	err := d.stagedService.Unstage(api.UnstageProductInput{ProductName: productName})
	if err != nil {
		return err
	}
	return nil
}

type dashboardService interface {
	GetInstallForm() (api.Form, error)
	GetRevertForm() (api.Form, error)
	PostInstallForm(api.PostFormInput) error
}

type RevertChangeService struct {
	dashboardService dashboardService
}

func NewRevertChangeService(ds dashboardService) RevertChangeService {
	return RevertChangeService{
		dashboardService: ds,
	}
}

func (r RevertChangeService) RevertChanges() error {
	revertChangeService := commands.NewRevertStagedChanges(r.dashboardService, logger.Info)
	err := revertChangeService.Execute([]string{""})
	if err != nil {
		return err
	}
	return nil
}

type PerformInstallService struct {
	authedClient network.OAuthClient
}

func (p PerformInstallService) PerformInstall(productName string) error {
	nullWriter := commands.NewLogWriter(ioutil.Discard)
	installationsService := api.NewInstallationsService(p.authedClient)

	applyChanges := commands.NewApplyChanges(installationsService, nullWriter, logger.Info, 10)
	err := applyChanges.Execute([]string{productName})
	if err != nil {
		return err
	}
	return nil
}