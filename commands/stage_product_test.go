package commands_test

import (
	"errors"
	"fmt"

	"github.com/rsamban/om/api"
	"github.com/rsamban/om/commands"
	"github.com/rsamban/om/commands/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StageProduct", func() {
	var (
		stagedProductsService    *fakes.ProductStager
		diagnosticService        *fakes.DiagnosticService
		availableProductsService *fakes.AvailableProductChecker
		logger                   *fakes.Logger
	)

	BeforeEach(func() {
		stagedProductsService = &fakes.ProductStager{}
		availableProductsService = &fakes.AvailableProductChecker{}
		diagnosticService = &fakes.DiagnosticService{}
		logger = &fakes.Logger{}
	})

	It("stages a product", func() {
		availableProductsService.CheckProductAvailabilityReturns(true, nil)

		command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)

		err := command.Execute([]string{
			"--product-name", "some-product",
			"--product-version", "some-version",
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(stagedProductsService.StageCallCount()).To(Equal(1))
		Expect(stagedProductsService.StageArgsForCall(0)).To(Equal(api.StageProductInput{
			ProductName:    "some-product",
			ProductVersion: "some-version",
		}))

		format, v := logger.PrintfArgsForCall(0)
		Expect(fmt.Sprintf(format, v...)).To(Equal("staging some-product some-version"))

		format, v = logger.PrintfArgsForCall(1)
		Expect(fmt.Sprintf(format, v...)).To(Equal("finished staging"))
	})

	Context("when the product version has already been staged", func() {
		It("no-ops and returns successfully", func() {
			availableProductsService.CheckProductAvailabilityReturns(true, nil)
			diagnosticService.ReportReturns(api.DiagnosticReport{
				StagedProducts: []api.DiagnosticProduct{
					{
						Name:    "some-product",
						Version: "some-version",
					},
				},
			}, nil)

			command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)

			err := command.Execute([]string{
				"--product-name", "some-product",
				"--product-version", "some-version",
			})
			Expect(err).NotTo(HaveOccurred())

			format, v := logger.PrintfArgsForCall(0)
			Expect(fmt.Sprintf(format, v...)).To(Equal("some-product some-version is already staged"))

			Expect(stagedProductsService.StageCallCount()).To(Equal(0))
		})
	})

	Context("failure cases", func() {
		Context("when an unknown flag is provided", func() {
			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)
				err := command.Execute([]string{"--badflag"})
				Expect(err).To(MatchError("could not parse stage-product flags: flag provided but not defined: -badflag"))
			})
		})

		Context("when the product-name flag is not provided", func() {
			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)
				err := command.Execute([]string{"--product-version", "1.0"})
				Expect(err).To(MatchError("error: product-name is missing. Please see usage for more information."))
			})
		})

		Context("when the product-version flag is not provided", func() {
			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)
				err := command.Execute([]string{"--product-name", "some-product"})
				Expect(err).To(MatchError("error: product-version is missing. Please see usage for more information."))
			})
		})

		Context("when the product is not available", func() {
			BeforeEach(func() {
				availableProductsService.CheckProductAvailabilityReturns(false, nil)
			})

			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)

				err := command.Execute([]string{
					"--product-name", "some-product",
					"--product-version", "some-version",
				})
				Expect(err).To(MatchError(ContainSubstring("failed to stage product: cannot find product")))
			})
		})

		Context("when the product availability cannot be determined", func() {
			BeforeEach(func() {
				availableProductsService.CheckProductAvailabilityReturns(false, errors.New("failed to check availability"))
			})

			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)

				err := command.Execute([]string{
					"--product-name", "some-product",
					"--product-version", "some-version",
				})
				Expect(err).To(MatchError(ContainSubstring("failed to stage product: cannot check availability")))
			})
		})

		Context("when the product cannot be staged", func() {
			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)
				availableProductsService.CheckProductAvailabilityReturns(true, nil)
				stagedProductsService.StageReturns(errors.New("some product error"))

				err := command.Execute([]string{"--product-name", "some-product", "--product-version", "some-version"})
				Expect(err).To(MatchError("failed to stage product: some product error"))
			})
		})

		Context("when the diagnostic report cannot be fetched", func() {
			It("returns an error", func() {
				command := commands.NewStageProduct(stagedProductsService, availableProductsService, diagnosticService, logger)
				availableProductsService.CheckProductAvailabilityReturns(true, nil)
				diagnosticService.ReportReturns(api.DiagnosticReport{}, errors.New("bad diagnostic report"))

				err := command.Execute([]string{"--product-name", "some-product", "--product-version", "some-version"})
				Expect(err).To(MatchError("failed to stage product: bad diagnostic report"))
			})
		})
	})

	Describe("Usage", func() {
		It("returns usage information for the command", func() {
			command := commands.NewStageProduct(nil, nil, nil, nil)
			Expect(command.Usage()).To(Equal(commands.Usage{
				Description:      "This command attempts to stage a product in the Ops Manager",
				ShortDescription: "stages a given product in the Ops Manager targeted",
				Flags:            command.Options,
			}))
		})
	})
})
