package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/koshatul/traefik-acme/traefik"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nolint: gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use:   "traefik-acme <domain>",
	Short: "Command to extract SSL certificates from traefik acme.json",
	Run:   mainCommand,
	Args:  cobra.MinimumNArgs(1),
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	cobra.OnInitialize(configInit)

	rootCmd.PersistentFlags().StringP("acme", "a", "/etc/traefik/acme.json", "Location of acme.json file")
	_ = viper.BindPFlag("acme", rootCmd.PersistentFlags().Lookup("acme"))
	_ = viper.BindEnv("acme", "ACME_FILE")

	rootCmd.PersistentFlags().StringP("cert", "c", "cert.pem", "Location to write out certificate")
	_ = viper.BindPFlag("cert", rootCmd.PersistentFlags().Lookup("cert"))
	_ = viper.BindEnv("cert", "CERT_FILE")

	rootCmd.PersistentFlags().StringP("key", "k", "key.pem", "Location to write out key file")
	_ = viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	_ = viper.BindEnv("key", "KEY_FILE")

	rootCmd.PersistentFlags().Bool("force", false, "Force writing to file even if not updated")
	_ = viper.BindPFlag("force", rootCmd.PersistentFlags().Lookup("force"))
	_ = viper.BindEnv("force", "FORCE_WRITE")

	rootCmd.PersistentFlags().Bool("exit-code", false, "Exit with exit-code 99 if files updated")
	_ = viper.BindPFlag("exit-code", rootCmd.PersistentFlags().Lookup("exit-code"))
	_ = viper.BindEnv("exit-code", "EXIT_CODE")

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")
}

func main() {
	_ = rootCmd.Execute()
}

//nolint: gocritic // ifElseChain doesn't seem to be idiomatic here.
func writeFile(filename string, data []byte, perm os.FileMode) (bool, error) {
	updated := false

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, just write it.
		logrus.WithField("filename", filename).Debugf("file not found, writing")

		updated = true
		err := ioutil.WriteFile(filename, data, perm)

		return updated, err
	} else if viper.GetBool("force") {
		// Don't care if it exists, just write it.
		logrus.WithField("filename", filename).Debugf("file found, but force enabled")

		updated = true
		err := ioutil.WriteFile(filename, data, perm)

		return updated, err
	} else {
		// File exists
		logrus.WithField("filename", filename).Debugf("file found")

		ld, err := ioutil.ReadFile(filename)
		if err != nil {
			return false, err
		}

		i := bytes.Compare(ld, data)
		if i == 0 {
			logrus.WithField("filename", filename).Debugf("file unchanged")
			return updated, nil
		}

		logrus.WithField("filename", filename).Debugf("file changed, writing")
		updated = true
		err = ioutil.WriteFile(filename, data, perm)

		return updated, err
	}
}

func mainCommand(cmd *cobra.Command, args []string) {
	domain := args[0]
	updated := false

	store, err := traefik.ReadFile(viper.GetString("acme"))
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	if cert := store.GetCertificateByName(domain); cert != nil {
		certUpdated, err := writeFile(viper.GetString("cert"), cert.Certificate, 0644)
		if err != nil {
			logrus.Errorf("unable to write certificate: %s", err.Error())
			os.Exit(1)
		}

		keyUpdated, err := writeFile(viper.GetString("key"), cert.Key, 0600)
		if err != nil {
			logrus.Errorf("unable to write key: %s", err.Error())
			os.Exit(1)
		}

		if certUpdated || keyUpdated {
			logrus.Printf("Successfully wrote %s certificate (%s) and key (%s)",
				domain,
				viper.GetString("cert"),
				viper.GetString("key"),
			)

			updated = true
		} else {
			logrus.Printf("Found %s, but certificate has not changed", domain)
		}
	} else {
		logrus.Printf("certificate not found for %s", domain)
		os.Exit(1)
	}

	if updated && viper.GetBool("exit-code") {
		os.Exit(99)
	}
}
