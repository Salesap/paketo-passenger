package passenger_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/passenger"
	"github.com/paketo-buildpacks/passenger/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir    string
		gemfileParser *fakes.Parser
		detect        packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(workingDir, "Gemfile"), []byte{}, 0644)
		Expect(err).NotTo(HaveOccurred())

		gemfileParser = &fakes.Parser{}

		detect = passenger.Detect(gemfileParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the Gemfile lists passenger and mri", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasPassenger = true
		})
		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "gems",
						Metadata: passenger.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "bundler",
						Metadata: passenger.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: "mri",
						Metadata: passenger.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			}))
		})
	})

	context("when the Gemfile does not list passenger", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasPassenger = false
		})

		it("detect should fail with error", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail))
		})
	})

	context("failure cases", func() {
		context("when the gemfile parser fails", func() {
			it.Before(func() {
				gemfileParser.ParseCall.Returns.Err = errors.New("some-error")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError("failed to parse Gemfile: some-error"))
			})
		})
	})
}
