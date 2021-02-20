package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/alex-held/devctl/internal/system"

	"github.com/alex-held/devctl/internal/sdkman"

	"github.com/alex-held/devctl/internal/cli"
)

var fmtFlag string

type OutputFormat string

const (
	Text  OutputFormat = "text"
	Table OutputFormat = "table"
)

// NewSdkManCommand creates the sdkman command
func NewSdkManCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "sdkman",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		NewSdkManListCommand(),
		NewSdkManVersionsCommand(),
		NewSdkManDefaultCommand(),
		NewSdkManDownloadCommand(),
	)
	return cmd
}

// NewSdkManCommand creates the sdkman command
func NewSdkManListCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list",
		Short: "lists all installable sdks",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := sdkman.NewSdkManClient()
			sdks, resp, err := client.ListSdks.ListAllSDK(ctx)
			if err != nil {
				cli.ExitWithError(1, err)
			}
			defer resp.Body.Close()

			formatFlag := cmd.Flag("format")
			format := formatFlag.Value.String()

			switch OutputFormat(format) {
			case Table:
				println(sdks.String())
				os.Exit(0)
			case Text:
				for _, sdk := range sdks {
					println(sdk)
				}
				os.Exit(0)
			default:
				println(sdks)
				os.Exit(0)
			}
		},
	}

	cmd.Flags().StringVarP(&fmtFlag, "format", "f", string(Table), "the output format of the cli app. -format=table")
	return cmd
}

// NewSdkManCommand creates the sdkman command
func NewSdkManDefaultCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:       "default",
		Short:     "Displays the latest (default) version of a given  sdk ",
		ValidArgs: strings.Split("ant,asciidoctorj,ballerina,bpipe,btrace,ceylon,concurnas,crash,cuba,cxf,doctoolchain,dotty,gaiden,glide,gradle,gradleprofiler,grails,groovy,groovyserv,http4k,infrastructor,java,jbake,jbang,karaf,kotlin,kscript,layrry,lazybones,leiningen,maven,micronaut,mulefd,mvnd,sbt,scala,spark,springboot,sshoogr,test,tomcat,vertx,visualvm", ","), //nolint: lll
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := sdkman.NewSdkManClient()
			defaultVersion, err := client.Version.Default(ctx, args[0])
			if err != nil {
				cli.ExitWithError(1, err)
			}

			fmt.Println(defaultVersion)
		},
	}
	return cmd
}

// NewSdkManCommand creates the sdkman command
func NewSdkManVersionsCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "versions",
		Short: "Displays the latest (default) version of a given  sdk ",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := sdkman.NewSdkManClient()
			defaultVersion, err := client.Version.All(ctx, args[0], system.MacOsx)
			if err != nil {
				cli.ExitWithError(1, err)
			}

			fmt.Println(defaultVersion)
		},
	}
	return cmd
}

// NewSdkManCommand creates the sdkman command
func NewSdkManDownloadCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		// 	Example: add [-F file | -D dir]... [-f format] profile
		Use:     `download`,
		Aliases: []string{"d", "dl"},
		Short:   "Downloads the sdk  (default) version of a given  sdk ",
		Long: `
		download [sdk] #downloads latest (default) version
		download [sdk] [version]
		download [sdk] [version] [system]`,
		Args: cobra.RangeArgs(1, 3),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			client := sdkman.NewSdkManClient()

			switch len(args) {
			case 0:
				cli.ExitWithError(1, os.ErrInvalid)
			case 1:
				sdk := args[0]
				version, err := client.Version.Default(ctx, sdk)
				cli.ExitWithError(1, err)
				dlPath := client.Download.Resolve()(sdk, version)

				dl, r, err := client.Download.DownloadSDK(ctx, dlPath, sdk, version, system.MacOsx)
				cli.ExitWithError(1, err)
				defer r.Body.Close()
				fmt.Printf("Downloaded sdk to path: %s", dl.Path)
				os.Exit(0)
			case 2: //nolint: gomnd
				sdk := args[0]
				version, err := semver.ParseTolerant(args[1])
				cli.ExitWithError(1, err)
				dlPath := client.Download.Resolve()(sdk, version.String())

				dl, r, err := client.Download.DownloadSDK(ctx, dlPath, sdk, version.String(), system.MacOsx)
				cli.ExitWithError(1, err)
				defer r.Body.Close()
				fmt.Printf("Downloaded sdk to path: %s", dl.Path)
				os.Exit(0)
			default:
				cli.ExitWithError(1, errors.New("sdkman download called with too many arguments"))
			}
		},
	}

	return cmd
}
