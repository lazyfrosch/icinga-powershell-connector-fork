package main

import (
	"crypto/tls"
	"fmt"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
)

const (
	DefaultAPI  = "https://localhost:5668"
	ProgramName = "icinga-powershell-connector"
)

type Config struct {
	API          string
	Command      string
	Arguments    map[string]interface{}
	CertName     string
	CAFile       string
	Insecure     bool
	Debug        bool
	PrintVersion bool
}

func ParseConfigFromFlags() (config *Config, err error) {
	config = &Config{}

	flag.StringVar(&config.Command, "command", "", "Command to be executed")
	flag.StringVar(&config.API, "api", DefaultAPI, "API Endpoint")
	flag.StringVar(&config.CertName, "cert-name", GetIcingaNodeName(), "Certificate Name to be expected")
	flag.StringVar(&config.CAFile, "ca-file", IcingaCAPath, "Icinga CA file to be loaded")
	flag.BoolVar(&config.Insecure, "insecure", false, "Ignore any certificate checks")
	flag.BoolVar(&config.Debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&config.PrintVersion, "version", false, "Print program version")

	flag.CommandLine.SortFlags = false
	flag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true

	flag.Parse()

	if config.PrintVersion {
		_, _ = fmt.Fprintln(os.Stdout, ProgramName+" "+buildVersion())
		_, _ = fmt.Fprintln(os.Stdout, License)

		os.Exit(0)
	}

	// Parse Powershell flags
	command, args := GetPowershellArgs(os.Args[1:])
	if command != "" {
		config.Command = command
	}

	config.Arguments = args

	if config.Command == "" {
		err = fmt.Errorf("no command found for Powershell execution")
		return
	}

	return
}

func (c Config) NewClient() *http.Client {
	tlsConfig := &tls.Config{
		RootCAs:            LoadIcingaCACert(c.CAFile),
		InsecureSkipVerify: c.Insecure, // nolint:gosec // intended configuration
		ServerName:         c.CertName,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &http.Client{Transport: transport}
}
