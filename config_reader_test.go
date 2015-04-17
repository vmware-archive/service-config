package service_config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/service-config"
)

var _ = Describe("ConfigReader", func() {
	const invalidYAML = `Count: INVALID`
	const simpleYAML = `Name: test-user`
	const nestedYAML = `---
Name: userName
Password: ppp
School: 
  Name: UB
  Location: Buffalo
`
	type ConfigSimple struct {
		Name string
	}
	type ConfigInvalid struct {
		Count int
	}
	type School struct {
		Name     string
		Location string
	}

	type ConfigNested struct {
		Name     string
		Password string
		School   School
	}

	Describe("Read", func() {

		It("unmarshal a config with one field", func() {
			reader := service_config.NewReader([]byte(simpleYAML))

			var simpleConfig ConfigSimple
			err := reader.Read(&simpleConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(simpleConfig).To(Equal(ConfigSimple{
				Name: "test-user",
			}))
		})

		It("unmarshal a config with nested fields", func() {
			reader := service_config.NewReader([]byte(nestedYAML))

			var nestedConfig ConfigNested
			err := reader.Read(&nestedConfig)
			Expect(err).NotTo(HaveOccurred())

			Expect(nestedConfig).To(Equal(ConfigNested{
				Name:     "userName",
				Password: "ppp",
				School: School{
					Name:     "UB",
					Location: "Buffalo",
				},
			}))
		})

		It("returns an error for unmarshalling a config without a valid YAML syntax", func() {
			reader := service_config.NewReader([]byte(invalidYAML))

			var invalidConfig ConfigInvalid
			err := reader.Read(&invalidConfig)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unmarshaling config"))
		})
	})
})
