package operations

import (
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"sync"
)

type ErrandState struct {
}

const (
	cfProduct     = "cf"
	rabbitProduct = "p-rabbitmq"
	mysqlProduct  = "p-mysql"
)

type errandDisabler interface {
	List(productID string) (api.ErrandsListOutput, error)
	SetState(productID, errandName string, postDeployState, preDeleteState interface{}) error
}

type DisableErrandsService struct {
	errandService errandDisabler
	stagedService stagedService
}

func NewDisableErrandService(es errandDisabler, ss stagedService) DisableErrandsService {
	return DisableErrandsService{
		errandService: es,
		stagedService: ss,
	}
}

func (e DisableErrandsService) DisableSystemErrands(wg *sync.WaitGroup) error {
	defer wg.Done()
	cfProductOutput, err := e.stagedService.Find(cfProduct)
	if err != nil {
		return err
	}
	rabbitProductOutput, err := e.stagedService.Find(rabbitProduct)
	if err != nil {
		return err
	}
	mysqlProductOutput, err := e.stagedService.Find(mysqlProduct)
	if err != nil {
		return err
	}
	cfErrands, err := e.errandService.List(cfProductOutput.Product.GUID)
	if err != nil {
		return err
	}
	rabbitErrands, err := e.errandService.List(rabbitProductOutput.Product.GUID)
	if err != nil {
		return err
	}
	mysqlErrands, err := e.errandService.List(mysqlProductOutput.Product.GUID)
	if err != nil {
		return err
	}
	disableErrandsForProduct(cfErrands, e, cfProductOutput.Product.GUID)
	disableErrandsForProduct(rabbitErrands, e, rabbitProductOutput.Product.GUID)
	disableErrandsForProduct(mysqlErrands, e, mysqlProductOutput.Product.GUID)
	return nil
}

func disableErrandsForProduct(productErrands api.ErrandsListOutput, e DisableErrandsService, productName string) {
	for _, productErrand := range productErrands.Errands {
		switch postDeployEnabled := productErrand.PostDeploy.(type) {
		case bool:
			if postDeployEnabled {
				setErrandState(e, productName, productErrand, false, nil)
			}
		case string:
			if postDeployEnabled == "default" || postDeployEnabled == "when-changed" {
				setErrandState(e, productName, productErrand, nil, false)
			}
		}

		switch preDeployEnabled := productErrand.PreDelete.(type) {
		case bool:
			if preDeployEnabled {
				setErrandState(e, productName, productErrand, nil, false)
			}
		case string:
			if preDeployEnabled == "default" || preDeployEnabled == "when-changed" {
				setErrandState(e, productName, productErrand, nil, false)
			}
		}
	}
}
func setErrandState(e DisableErrandsService, productName string, productErrand api.Errand, postDeploy, preDelete interface{}) {
	err := e.errandService.SetState(productName, productErrand.Name, postDeploy, preDelete)
	if err != nil {
		logger.Error.Printf("Error disabling errand %s for product %s: %s", productErrand.Name, productName, err)
	} else {
		logger.Info.Printf("disabling errand %s for product %s", productErrand.Name, productName)
	}
}
