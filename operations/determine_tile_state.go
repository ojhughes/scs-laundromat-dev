package operations

import (
	"errors"
	"fmt"
	"github.com/pivotal-cf/scs-laundromat-dev/logger"
	"github.com/pivotal-cf/om/api"
	"os"
	"strings"
)

const (
	StagedButNotDeployed TileState = 1 + iota
	StagedAndDeployed
	PendingUpdate
	PendingDeletion
	NotStaged
)
const (
	Success PostInstallState = 1 + iota
	ErrandFailed
)

type TileState int
type PostInstallState int

type stagedService interface {
	StagedProducts() (api.StagedProductsOutput, error)
	Unstage(input api.UnstageProductInput) error
}

type deployedService interface {
	DeployedProducts() ([]api.DeployedProductOutput, error)
}

type pendingChangeService interface {
	List() (api.PendingChangesOutput, error)
}

type GetStateInput struct {
	ProductName string
}

type DetermineTileStateService struct {
	stageService         stagedService
	deployedService      deployedService
	pendingChangeService pendingChangeService
}

func NewDetermineTileStateService(ss stagedService, ds deployedService, ps pendingChangeService) DetermineTileStateService {
	return DetermineTileStateService{
		stageService:         ss,
		deployedService:      ds,
		pendingChangeService: ps,
	}
}

func (d DetermineTileStateService) PreInstallState(input GetStateInput) (TileState, error) {
	isScsTilePendingDeletion := false
	isScsTilePendingUpdate := false

	stagedProducts, err := d.stageService.StagedProducts()
	if err != nil {
		return 0, err
	}
	deployedProducts, err := d.deployedService.DeployedProducts()
	if err != nil {
		return 0, err
	}

	isScsTileDeployed := isScsTileDeployed(deployedProducts, input)
	isScsTileStaged := isScsTileStaged(stagedProducts, input)

	if isScsTileStaged && !isScsTileDeployed {
		logger.Info.Println("SCS tile is staged but has not been installed")
		return StagedButNotDeployed, nil
	}

	if !isScsTileStaged && !isScsTileDeployed {
		logger.Info.Printf("SCS tile is not staged or deployed")
		return NotStaged, nil
	}

	if isScsTileDeployed && isScsTileStaged {
		pendingChanges, err := d.pendingChangeService.List()
		if err != nil {
			return 0, err
		}
		for _, change := range pendingChanges.ChangeList {
			if strings.HasPrefix(change.Product, input.ProductName) && change.Action == "delete" {
				isScsTilePendingDeletion = true

			} else if strings.HasPrefix(change.Product, input.ProductName) && change.Action == "update" {
				isScsTilePendingUpdate = true
			}
		}
		if isScsTilePendingDeletion {
			logger.Info.Println("SCS tile is already pending deletion update from a previous action")
			return PendingDeletion, nil
		} else if isScsTilePendingUpdate {
			logger.Info.Println("SCS tile is pending an update from a previous action")
			return PendingUpdate, nil
		} else {
			logger.Info.Println("SCS tile is staged and installed with no pending changed")
			return StagedAndDeployed, nil
		}
	}
	return 0, errors.New("unable to determine tile state")
}

func isScsTileStaged(stagedProducts api.StagedProductsOutput, input GetStateInput) bool {
	isStaged := false
	for _, stagedProduct := range stagedProducts.Products {
		if input.ProductName == stagedProduct.Type {
			fmt.Println("SCS is staged")
			isStaged = true
		}
	}
	return isStaged
}
func isScsTileDeployed(deployedProducts []api.DeployedProductOutput, input GetStateInput) bool {
	isDeployed := false
	for _, deployedProduct := range deployedProducts {
		if input.ProductName == deployedProduct.Type {
			isDeployed = true
			fmt.Println("SCS is deployed")
		}
	}
	return isDeployed
}
