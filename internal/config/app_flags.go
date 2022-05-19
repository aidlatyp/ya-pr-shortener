package config

import "github.com/spf13/pflag"

type AppFlags struct {
	addr     *string
	baseURL  *string
	fileName *string
}

// Addr and other methods to get unexported fields
// if this structure become too large and will be moved to separate package
// Also these methods satisfy configuration interfaces
func (p *AppFlags) Addr() string {
	return *p.addr
}

func (p *AppFlags) BaseURL() string {
	return *p.baseURL
}

func (p *AppFlags) Filename() string {
	return *p.fileName
}

func ParseFlags() AppFlags {

	parsed := AppFlags{}
	parsed.addr = pflag.StringP("a", "a", "", "Host IP address")
	parsed.baseURL = pflag.StringP("b", "b", "", "Base URL")
	parsed.fileName = pflag.StringP("f", "f", "", "Filename to store URLs")

	pflag.Parse()

	return parsed
}
