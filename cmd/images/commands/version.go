package commands

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Global variables set at build time.
	//
	// You typically build you app with flags such as:
	//
	//   go build \
	//     --ldflags "-s -w \
	//     -X 'github.com/fredbi/demo-api/cmd/images/commands.Version=$(git describe --tags) \
	//     -X 'github.com/fredbi/demo-api/cmd/images/commands.BuildDate=$(date -u -R) \
	//     -X 'github.com/fredbi/demo-api/cmd/images/commands.GitCommit=$(git rev-parse HEAD) \
	//     -X 'github.com/fredbi/demo-api/cmd/images/commands.GitState=$(git diff --quiet || echo dirty)' \
	//   " ./cmd/images

	// Version of this build, e.g. git tag
	Version string
	// BuildDate is a date as string (unconstrained)
	BuildDate string
	// GitCommit is the git hash
	GitCommit string
	// GitState reflects the current build state in git (i.e. clean / dirty)
	GitState string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the version of this API",
	Long: `Prints the version of this API. It includes the following components:

	* Semver (output of git describe --tags)
	* Build Date (date at which the binary was built)
	* Git Commit (the git commit hash this binary was built from)
	* Git State (the current git repo is dirty there are uncommitted changes used by the build)
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(NewVersionInfo().String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// NewVersionInfo build a new version object from global variables set a build time
func NewVersionInfo() VersionInfo {
	ver := VersionInfo{
		Version:   "dev",
		BuildDate: BuildDate,
		GitCommit: GitCommit,
		GitState:  "",
	}
	if Version != "" {
		ver.Version = Version
		ver.GitState = "clean"
	}
	if GitState != "" {
		ver.GitState = GitState
	}
	return ver
}

// VersionInfo completely describes a build version
type VersionInfo struct {
	Version   string `json:"version,omitempty"`
	BuildDate string `json:"buildDate,omitempty"`
	GitCommit string `json:"gitCommit,omitempty"`
	GitState  string `json:"gitState,omitempty"`
}

func (v VersionInfo) String() string {
	var buf bytes.Buffer
	buf.WriteString("Version: ")
	buf.WriteString(v.Version)
	buf.WriteString("\n")
	buf.WriteString("Build date: ")
	buf.WriteString(v.BuildDate)
	buf.WriteString("\n")
	buf.WriteString("Commit: ")
	buf.WriteString(v.GitCommit)
	buf.WriteString("\n")
	buf.WriteString("Working tree: ")
	buf.WriteString(v.GitState)
	buf.WriteString("\n")

	return buf.String()
}

/* TODO(fred)
func VersionHandler(log middleware.Logger, info VersionInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		enc := json.NewEncoder(w)
		if err := enc.Encode(info); err != nil {
			log.Printf("failed to write version response: %v", err)
		}
	}
}
*/
