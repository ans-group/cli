module github.com/ukfast/cli

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/golang/mock v1.2.0
	github.com/google/uuid v1.1.1
	github.com/guptarohit/asciigraph v0.4.2-0.20191006150553-f9506970428c
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/rhysd/go-github-selfupdate v1.1.0
	github.com/ryanuber/go-glob v1.0.0
	github.com/spf13/afero v1.2.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.6.1
	github.com/ukfast/sdk-go v1.3.25
	golang.org/x/sys v0.0.0-20190507160741-ecd444e8653b // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1
	k8s.io/client-go v11.0.0+incompatible
)

// replace github.com/ukfast/sdk-go => ../sdk-go
