package root_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Root", func() {
	Context("--help", func() {
		It("should output help text", func() {
			output, err := cli.Help()
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(ContainSubstring("Launch and manage applications on Amazon ECS and AWS Fargate."))
			Expect(output).To(ContainSubstring("Getting Started"))
			Expect(output).To(ContainSubstring("Develop"))
			Expect(output).To(ContainSubstring("Release"))
			Expect(output).To(ContainSubstring("Settings"))
		})
	})

	Context("--version", func() {
		It("should output a valid semantic version", func() {
			output, err := cli.Version()
			Expect(err).NotTo(HaveOccurred())
			// Versions look like ecs-preview version: v0.0.4-34-g133b977
			// the extra bit at the end is if the build isn't a tagged release.
			Expect(output).To(MatchRegexp(`ecs-preview version: v\d*\.\d*\.\d*.*`))
		})
	})
})
