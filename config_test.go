package service_config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

const configJSON = `{
    "Name": "Enterprise",
    "ID": 1701,
    "Crew": {
        "Officers": [
            {"Name": "Kirk", "Role": "Commanding Officer"},
            {"Name": "Kirk", "Role": "First Officer/Science Officer"},
            {"Name": "McCoy", "Role": "Chief Medical Officer"}
        ],
        "Passengers": [
            {"Name": "Sarek", "Title": "Federation Ambassador"}
        ]
    }
}`

const configJSONAlt = `{
    "Name": "Defiant",
    "ID": 74205,
    "Crew": {
        "Officers": [
            {"Name": "Sisko", "Role": "Commanding Officer"},
            {"Name": "Worf", "Role": "Strategic Operations Officer"},
        ]
    }
}`

const configStructString = `main.ShipConfig{Name:"Enterprise", ID:1701, Crew:main.Crew{Officers:[]main.Officer{main.Officer{Name:"Kirk", Role:"Commanding Officer"}, main.Officer{Name:"Kirk", Role:"First Officer/Science Officer"}, main.Officer{Name:"McCoy", Role:"Chief Medical Officer"}}, Passengers:[]main.Passenger{main.Passenger{Name:"Sarek", Title:"Federation Ambassador"}}}}`

var _ = Describe("ServiceConfig", func() {
    var whitespacePattern = regexp.MustCompile("\\s+")
	var command *exec.Cmd

	Context("When a config flag is passed", func() {
		BeforeEach(func() {
			configString := whitespacePattern.ReplaceAllString(configJSON, " ")
			command = exec.Command(
				binaryPath,
				fmt.Sprintf("-config=%s", configString),
			)
		})

		It("Reads the config, from the flag string", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			session.Wait(30 * time.Second)
			Expect(session).To(gexec.Exit(0))

			Expect(string(session.Out.Contents())).To(ContainSubstring("Config: %s", configStructString))
		})

        Context("When the CONFIG env var is ALSO set", func() {
            BeforeEach(func() {
                configString := whitespacePattern.ReplaceAllString(configJSONAlt, " ")

                command.Env = []string{
                    fmt.Sprintf("CONFIG=%s", configString),
                }
            })

            It("Reads the config, from the flag string", func() {
                session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
                Expect(err).ToNot(HaveOccurred())

                session.Wait(30 * time.Second)
                Expect(session).To(gexec.Exit(0))

                Expect(string(session.Out.Contents())).To(ContainSubstring("Config: %s", configStructString))
            })
        })
	})

	Context("When a configPath flag is passed", func() {
		BeforeEach(func() {
			configPath := filepath.Join(tempDir, "flag-config.json")

			err := ioutil.WriteFile(configPath, []byte(configJSON), os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			command = exec.Command(
				binaryPath,
				fmt.Sprintf("-configPath=%s", configPath),
			)
		})

		It("Reads the config, from the file path specified by flag", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			session.Wait(30 * time.Second)
			Expect(session).To(gexec.Exit(0))

			Expect(string(session.Out.Contents())).To(ContainSubstring("Config: %s", configStructString))
		})

        Context("When a CONFIG_PATH env var is ALSO set", func() {
            BeforeEach(func() {
                configPath := filepath.Join(tempDir, "flag-env-var-config.json")

                err := ioutil.WriteFile(configPath, []byte(configJSONAlt), os.ModePerm)
                Expect(err).ToNot(HaveOccurred())

                command.Env = []string{
                    fmt.Sprintf("CONFIG_PATH=%s", configPath),
                }
            })

            It("Reads the config, from the file path specified by flag", func() {
                session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
                Expect(err).ToNot(HaveOccurred())

                session.Wait(30 * time.Second)
                Expect(session).To(gexec.Exit(0))

                Expect(string(session.Out.Contents())).To(ContainSubstring("Config: %s", configStructString))
            })
        })
	})

	Context("When a CONFIG env var is set", func() {
		BeforeEach(func() {
			configString := whitespacePattern.ReplaceAllString(configJSON, " ")

			command = exec.Command(binaryPath)
			command.Env = []string{
				fmt.Sprintf("CONFIG=%s", configString),
			}
		})

		It("Reads the config, from the env var string", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			session.Wait(30 * time.Second)
			Expect(session).To(gexec.Exit(0))

			Expect(string(session.Out.Contents())).To(ContainSubstring("Config: %s", configStructString))
		})
	})

	Context("When a CONFIG_PATH env var is set", func() {
		BeforeEach(func() {
			configPath := filepath.Join(tempDir, "env-var-config.json")

			err := ioutil.WriteFile(configPath, []byte(configJSON), os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			command = exec.Command(binaryPath)
			command.Env = []string{
				fmt.Sprintf("CONFIG_PATH=%s", configPath),
			}
		})

		It("Reads the config, from the file path specified by env var", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			session.Wait(30 * time.Second)
			Expect(session).To(gexec.Exit(0))

			Expect(string(session.Out.Contents())).To(ContainSubstring("Config: %s", configStructString))
		})
	})
})
