package commands

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf/jhanda/commands"
	"github.com/pivotal-cf/jhanda/flags"
	"github.com/pivotal-cf/om/api"
)

type StageProduct struct {
	logger                   logger
	stagedProductsService    productStager
	deployedProductsService  deployedProductsLister
	availableProductsService availableProductChecker
	diagnosticService        diagnosticService
	Options                  struct {
		Product string `short:"p"  long:"product-name"  description:"name of product"`
		Version string `short:"v"  long:"product-version"  description:"version of product"`
	}
}

//go:generate counterfeiter -o ./fakes/product_stager.go --fake-name ProductStager . productStager
type productStager interface {
	Stage(api.StageProductInput, string) error
}

//go:generate counterfeiter -o ./fakes/deployed_products_lister.go --fake-name DeployedProductsLister . deployedProductsLister
type deployedProductsLister interface {
	DeployedProducts() ([]api.DeployedProductOutput, error)
}

//go:generate counterfeiter -o ./fakes/available_product_checker.go --fake-name AvailableProductChecker . availableProductChecker
type availableProductChecker interface {
	CheckProductAvailability(productName string, productVersion string) (bool, error)
}

func NewStageProduct(productStager productStager, deployedProductsService deployedProductsLister, availableProductChecker availableProductChecker, diagnosticService diagnosticService, logger logger) StageProduct {
	return StageProduct{
		logger:                   logger,
		stagedProductsService:    productStager,
		deployedProductsService:  deployedProductsService,
		availableProductsService: availableProductChecker,
		diagnosticService:        diagnosticService,
	}
}

func (sp StageProduct) Execute(args []string) error {
	_, err := flags.Parse(&sp.Options, args)
	if err != nil {
		return fmt.Errorf("could not parse stage-product flags: %s", err)
	}

	if sp.Options.Product == "" {
		return errors.New("error: product-name is missing. Please see usage for more information.")
	}

	if sp.Options.Version == "" {
		return errors.New("error: product-version is missing. Please see usage for more information.")
	}

	diagnosticReport, err := sp.diagnosticService.Report()
	if err != nil {
		return fmt.Errorf("failed to stage product: %s", err)
	}

	deployedProductGUID := ""
	deployedProducts, err := sp.deployedProductsService.DeployedProducts()
	for _, deployedProduct := range deployedProducts {
		if deployedProduct.Type == sp.Options.Product {
			deployedProductGUID = deployedProduct.GUID
			break
		}
	}
	if err != nil {
		return fmt.Errorf("failed to stage product: %s", err)
	}

	for _, stagedProduct := range diagnosticReport.StagedProducts {
		if stagedProduct.Name == sp.Options.Product && stagedProduct.Version == sp.Options.Version {
			sp.logger.Printf("%s %s is already staged", sp.Options.Product, sp.Options.Version)
			return nil
		}
	}

	available, err := sp.availableProductsService.CheckProductAvailability(sp.Options.Product, sp.Options.Version)
	if err != nil {
		return fmt.Errorf("failed to stage product: cannot check availability of product %s %s", sp.Options.Product, sp.Options.Version)
	}

	if !available {
		return fmt.Errorf("failed to stage product: cannot find product %s %s", sp.Options.Product, sp.Options.Version)
	}

	sp.logger.Printf("staging %s %s", sp.Options.Product, sp.Options.Version)

	err = sp.stagedProductsService.Stage(api.StageProductInput{
		ProductName:    sp.Options.Product,
		ProductVersion: sp.Options.Version,
	}, deployedProductGUID)
	if err != nil {
		return fmt.Errorf("failed to stage product: %s", err)
	}

	sp.logger.Printf("finished staging")

	return nil
}

func (sp StageProduct) Usage() commands.Usage {
	return commands.Usage{
		Description:      "This command attempts to stage a product in the Ops Manager",
		ShortDescription: "stages a given product in the Ops Manager targeted",
		Flags:            sp.Options,
	}
}
