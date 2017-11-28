package operations

import (
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"time"
)

const (
	uaaAdminCredentialReference = ".uaa.admin_credentials"
	systemOrgName               = "system"
	serviceInstancesOrgName     = "p-spring-cloud-services"
	systemSpaceName             = "p-spring-cloud-services"
	serviceInstancesSpaceName   = "instances"
)

type RemoveIncompleteTileService struct {
	stagedService stagedService
}

func NewRemoveIncompleteTileService(ss stagedService) RemoveIncompleteTileService {
	return RemoveIncompleteTileService{
		stagedService: ss,
	}
}

func (d RemoveIncompleteTileService) RemoveInstall(productName string) error {
	logger.Info.Println("Unstaging tile")
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
	logger.Info.Println("Reverting pending changes")
	revertChangeService := commands.NewRevertStagedChanges(r.dashboardService, logger.Info)
	err := revertChangeService.Execute([]string{""})
	if err != nil {
		return err
	}
	return nil
}

type installService interface {
	Trigger(bool, bool) (api.InstallationsServiceOutput, error)
	Status(id int) (api.InstallationsServiceOutput, error)
}

type PerformInstallService struct {
	installationsService installService
}

func NewPerformInstallService(is installService) PerformInstallService {
	return PerformInstallService{
		installationsService: is,
	}
}

func (p PerformInstallService) PerformInstall(productName string) (api.InstallationsServiceOutput, error) {
	logger.Info.Println("Applying changes to delete tile")
	emptyOutput := api.InstallationsServiceOutput{}
	installation, err := p.installationsService.Trigger(true, true)
	if err != nil {
		return emptyOutput, err
	}
	for {
		current, err := p.installationsService.Status(installation.ID)
		if err != nil {
			return emptyOutput, fmt.Errorf("installation failed to get status: %s", err)
		}

		if current.Status == api.StatusSucceeded {
			return emptyOutput, nil
		} else if current.Status == api.StatusFailed {
			return emptyOutput, errors.New("installation was unsuccessful")
		}

		time.Sleep(time.Duration(10) * time.Second)
	}
}

type credentialsService interface {
	Fetch(deployedProductGUID, credentialReference string) (api.CredentialOutput, error)
}

type ForceUninstallService struct {
	cfClient CfClientService
}

type CfClientFactory interface {
	NewClient(config *cfclient.Config) (CfClientService, error)
}

type CfClientService interface {
	GetOrgByName(name string) (cfclient.Org, error)
	GetSpaceByName(spaceName string, orgGuid string) (space cfclient.Space, err error)
	DeleteSpace(guid string, recursive, async bool) error
}

func NewCfClientFactory() CfClientFactory {
	return &cfClientFactory{}
}

type cfClientFactory struct{}

func (self *cfClientFactory) NewClient(config *cfclient.Config) (CfClientService, error) {
	return cfclient.NewClient(config)
}

func NewForceUninstallService(cc CfClientService) ForceUninstallService {
	return ForceUninstallService{
		cfClient: cc,
	}
}

func (f ForceUninstallService) ForceUninstall(productName string) error {
	logger.Info.Println("Forcefully uninstalling SCS by removing tile orgs and spaces")
	systemOrgGuid, err := f.cfClient.GetOrgByName(systemOrgName)
	if err != nil {
		return err
	}
	systemSpaceGuid, err := f.cfClient.GetSpaceByName(systemSpaceName, systemOrgGuid.Guid)
	if err != nil {
		return err
	}
	serviceInstancesOrgGuid, err := f.cfClient.GetOrgByName(serviceInstancesOrgName)
	if err != nil {
		return err
	}
	serviceInstancesSpaceGuid, err := f.cfClient.GetSpaceByName(serviceInstancesSpaceName, serviceInstancesOrgGuid.Guid)
	if err != nil {
		return err
	}
	err = f.cfClient.DeleteSpace(systemSpaceGuid.Guid, true, false)
	if err != nil {
		return err
	}
	err = f.cfClient.DeleteSpace(serviceInstancesSpaceGuid.Guid, true, false)
	if err != nil {
		return err
	}
	return nil
}

func NewExtractCredentialService(cs credentialsService, ds deployedService) ExtractCredentialService {
	return ExtractCredentialService{
		credentialsService: cs,
		deployedService:    ds,
	}
}

type ExtractCredentialService struct {
	credentialsService credentialsService
	deployedService    deployedService
}

func (ecs ExtractCredentialService) ExtractCfPassword(productName string) (map[string]string, error) {
	emptyMap := make(map[string]string)
	deployedProductGUID := ""
	deployedProducts, err := ecs.deployedService.DeployedProducts()
	if err != nil {
		return emptyMap, fmt.Errorf("failed to fetch credential: %s", err)
	}
	for _, deployedProduct := range deployedProducts {
		if deployedProduct.Type == productName {
			deployedProductGUID = deployedProduct.GUID
			break
		}
	}

	output, err := ecs.credentialsService.Fetch(deployedProductGUID, uaaAdminCredentialReference)
	if err != nil {
		return emptyMap, fmt.Errorf("failed to fetch credential for %q: %s", uaaAdminCredentialReference, err)
	}
	return output.Credential.Value, nil
}
