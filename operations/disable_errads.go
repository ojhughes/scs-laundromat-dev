package operations

import (
	"sync"
	"github.com/pivotal-cf/om/api"
)


type errandDisabler interface {
	List(productID string) (api.ErrandsListOutput, error)
	SetState(productID, errandName string, postDeployState, preDeleteState interface{}) error
}
type DisableErrandsService struct {
	errandService errandDisabler
}

func NewDisableErrandService(es errandDisabler) DisableErrandsService {
	return DisableErrandsService{
		errandService: es,
	}
}

func (e DisableErrandsService) DisableSystemErrands(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
}
