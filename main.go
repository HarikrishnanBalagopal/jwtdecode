package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	VERSION         = "v1.0.3"
	GIT_TREE_STATE  = "unknown"
	GIT_COMMIT_HASH = "unknown"
)

var (
	GO_VERSION = runtime.Version()
)

func handleVersion(long bool) {
	fmt.Println(VERSION)
	if long {
		fmt.Println("gitCommit:", GIT_COMMIT_HASH)
		fmt.Println("gitTreeState:", GIT_TREE_STATE)
		fmt.Println("goVersion:", GO_VERSION)
	}
}

func getVersionCommand() *cobra.Command {
	long := false
	cmd := &cobra.Command{
		Use: "version",
		Run: func(*cobra.Command, []string) {
			handleVersion(long)
		},
	}
	cmd.Flags().BoolVarP(&long, "long", "l", false, "Print more detailed version info. Git commit, Git tree state, etc.")
	return cmd
}

func handleDecode(jwt string) (header string, payload string, err error) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("expected there to be 3 parts. actual len %d . actual: %s", len(parts), jwt)
	}
	for i, part := range parts[:2] {
		x, err := base64.RawURLEncoding.DecodeString(part)
		if err != nil {
			return "", "", err
		}
		parts[i] = string(x)
	}
	return parts[0], parts[1], nil
}

func handleCLI(jwt string) {
	if jwt == "-" {
		jwtBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			logrus.Fatalf("failed to read from stdin. Error: %q", err)
		}
		jwt = strings.TrimSpace(string(jwtBytes))
	}
	header, payload, err := handleDecode(jwt)
	if err != nil {
		logrus.Fatalf("failed to decode the jwt. Error: %q", err)
	}
	logrus.Debugf("header: %s", header)
	fmt.Print(payload)
}

func main() {
	verbose := false
	rootCmd := &cobra.Command{
		Use:   "jwtdecode ${MY_JWT}",
		Short: "jwtdecode ${MY_JWT} or echo ${MY_JWT} | jwtdecode -",
		Long: `jwtdecode helps decode JSON Web Tokens (JWTs).
A JWT can be encrypted (JWE) or just signed (JWS).
Support is only there for JWSs.
This program expects the JWT to be in compact serialization format:

<base64url encoded header>.<base64url encoded payload>.<base64url encoded signature>
`,
		Args: cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if verbose {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}
			handleCLI(args[0])
		},
	}
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print more logging messages. Useful for debugging.")
	rootCmd.AddCommand(getVersionCommand())
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error: %q", err)
	}
}
