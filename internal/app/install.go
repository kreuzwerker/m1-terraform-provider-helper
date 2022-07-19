package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	goversion "github.com/hashicorp/go-version"
)

const requestTimeoutSeconds int = 2

type Provider struct {
	Repo        string `json:"source"`
	Description string `json:"description"`
}

type BuildCommandInformation struct {
	command         string
	startingVersion *goversion.Version
}

type TerraformVersion struct {
	Version string `json:"terraform_version"`
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func executeBashCommand(command string, baseDir string) {
	shExecutable, _ := exec.LookPath("sh")

	cmd := &exec.Cmd{
		Path:   shExecutable,
		Args:   []string{shExecutable, "-c", command},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    baseDir,
	}

	if err := cmd.Run(); err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			if e.ExitCode() != 0 {
				log.Fatalf("Bash execution did not run successfully: %s", err)
			}
		}
	}
}

func executeBashCommandAndReturnOutput(command string, baseDir string) string {
	shExecutable, _ := exec.LookPath("sh")

	cmd := &exec.Cmd{
		Path:   shExecutable,
		Args:   []string{shExecutable, "-c", command},
		Stdout: nil,
		Stderr: os.Stderr,
		Dir:    baseDir,
	}

	output, err := cmd.Output()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			if e.ExitCode() != 0 {
				log.Fatalf("Bash execution did not run successfully: %s", err)
			}
		}
	}

	return string(output)
}

func getProviderData(providerName string) (Provider, error) {
	url := "https://registry.terraform.io/v1/providers/" + providerName

	client := &http.Client{Timeout: time.Second * time.Duration(float64(requestTimeoutSeconds))}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return Provider{}, fmt.Errorf("request error %w", err)
	}

	res, err := client.Do(req)

	if err != nil {
		if os.IsTimeout(err) {
			return Provider{}, fmt.Errorf("timeout error while trying to get provider data from "+url+": %w", err)
		}

		return Provider{}, fmt.Errorf("response error %w", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return Provider{}, fmt.Errorf("body reading error %w", err)
	}

	var data Provider

	err = json.Unmarshal(body, &data)
	if err != nil {
		return Provider{}, fmt.Errorf("could not parse JSON %w", err)
	}

	return data, nil
}

func cloneRepo(gitURL string, fullPath string) {
	_, err := git.PlainClone(fullPath, false, &git.CloneOptions{
		URL:      gitURL,
		Progress: os.Stdout,
	})

	CheckIfError(err)
}

// if repo is not check out yet
//		- clone to cli dir
// if already exists: dont clone, simply cd
// on both casees: checkout version
// return path to dir.
func checkoutSourceCode(baseDir string, gitURL string, version string) string {
	var r *git.Repository

	repoDir := extractRepoNameFromURL(gitURL)
	fullPath := baseDir + "/" + repoDir

	if !isDirExistent(fullPath) {
		cloneRepo(gitURL, fullPath)
	}

	r, err := git.PlainOpen(fullPath)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// Clean the repository
	executeBashCommand("git reset --hard && git clean -d -f -q", fullPath)
	log.Println("Pulling newest changes from " + gitURL)
	executeBashCommand("git remote show origin | sed -n '/HEAD branch/s/.*: //p'| xargs git checkout && git pull", fullPath)

	if len(version) > 0 {
		log.Println("version: " + version)
		ref, _ := r.ResolveRevision(plumbing.Revision(version))
		err = w.Checkout(&git.CheckoutOptions{
			Hash: *ref,
		})
		CheckIfError(err)
	} else {
		log.Println("No version specified, staying on latest commit")
	}

	return repoDir
}

func extractRepoNameFromURL(url string) string {
	parts := strings.Split(url, "/")

	return parts[len(parts)-1]
}

func normalizeSemver(version string) string {
	if strings.HasPrefix(version, "v") {
		return version[1:]
	}

	return version
}

func createBuildCommand(providerName string, version string, goPath string) string {
	parsedVersion, err := goversion.NewVersion(version)
	CheckIfError(err)

	v0, _ := goversion.NewVersion("0")
	v1, _ := goversion.NewVersion("1")
	v2, _ := goversion.NewVersion("2")
	v3, _ := goversion.NewVersion("3")
	v4, _ := goversion.NewVersion("4")

	buildCommands := make(map[string][]BuildCommandInformation)
	buildCommands["default"] = []BuildCommandInformation{{command: "make build", startingVersion: v0}}
	buildCommands["hashicorp/helm"] = []BuildCommandInformation{{command: "make build && cp terraform-provider-helm " + goPath + "/bin/" + "terraform-provider-helm", startingVersion: v0}}
	buildCommands["hashicorp/google"] = []BuildCommandInformation{{command: "gofmt -s -w ./tools.go  && make build", startingVersion: v0}}
	buildCommands["hashicorp/aws"] = []BuildCommandInformation{
		{command: "make tools && make fmt && gofmt -s -w ./tools.go && make build", startingVersion: v0},
		{command: "go mod init && go mod vendor && make fmt && make build", startingVersion: v1},
		{command: "make tools && make fmt && gofmt -s -w ./tools.go && make build", startingVersion: v2},
		{command: "cd tools && go get -d github.com/pavius/impi/cmd/impi && cd .. && make tools && make build", startingVersion: v3},
		{command: "make tools && make build", startingVersion: v4},
	}

	buildCommandMap, exists := buildCommands[providerName]

	if exists {
		var foundBuildCommand string

		for _, v := range buildCommandMap {
			if parsedVersion.GreaterThanOrEqual(v.startingVersion) {
				foundBuildCommand = v.command
			}
		}

		return foundBuildCommand
	}

	return buildCommands["default"][0].command
}

func (a *App) buildProvider(dir string, providerName string, version string, customBuildCommand string) {
	var buildCommand string

	if len(customBuildCommand) > 0 {
		fmt.Fprintf(os.Stdout, "Using custom build command: \"%s\"\n", customBuildCommand)
		buildCommand = customBuildCommand
	} else {
		buildCommand = createBuildCommand(providerName, version, a.Config.GoPath)
	}
	// #nosec G204
	executeBashCommand(buildCommand, a.Config.ProvidersCacheDir+"/"+dir)
}

func (a *App) moveBinaryToCorrectLocation(providerName string, version string, executableName string) {
	if len(version) == 0 {
		version = "master"
	} else {
		version = normalizeSemver(version)
	}

	newPath := a.createDestinationAndReturnExecutablePath(providerName, version, executableName)

	pathOfExecutable := a.Config.GoPath + "/bin/" + executableName
	log.Print("Move from " + pathOfExecutable + " to " + newPath)
	err := os.Rename(pathOfExecutable, newPath)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) createDestinationAndReturnExecutablePath(providerName string, version string, executableName string) string {
	oldTfVersion, _ := goversion.NewVersion("1.2.0")
	currentTfVersion := getTerraformVersion()

	var newPath string

	if currentTfVersion.GreaterThan(oldTfVersion) {
		filePath := a.Config.TerraformPluginDir + "/registry.terraform.io/" + providerName + "/" + version + "/darwin_arm64"
		createDirIfNotExists(filePath)

		fmt.Fprintf(os.Stdout, "GOPATH: %s\n", a.Config.GoPath)

		newPath = filePath + "/" + executableName + "_" + version + "_x5"
	} else {
		// before 0.12.31 it is: ~/.terraform.d/plugins/darwin_arm64/terraform-provider-template_v2.2.0
		filePath := a.Config.TerraformPluginDir + "/darwin_arm64"
		createDirIfNotExists(filePath)

		fmt.Fprintf(os.Stdout, "GOPATH: %s\n", a.Config.GoPath)
		newPath = filePath + "/" + executableName + "_v" + version
	}

	return newPath
}

func getTerraformVersion() *goversion.Version {
	versionRaw := executeBashCommandAndReturnOutput("terraform version -json", "./")

	var versionObj TerraformVersion

	err := json.Unmarshal([]byte(versionRaw), &versionObj)
	CheckIfError(err)
	parsedVersion, err := goversion.NewVersion(versionObj.Version)
	CheckIfError(err)

	return parsedVersion
}

func (a *App) Install(providerName string, version string, customBuildCommand string) bool {
	providerData, err := getProviderData(providerName)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Error while trying to get provider data from terraform registry: %v", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Repo: %s\n", providerData.Repo)

	gitRepo := providerData.Repo
	fmt.Fprintf(os.Stdout, "GitRepo: %s\n", gitRepo)

	sourceCodeDir := checkoutSourceCode(a.Config.ProvidersCacheDir, gitRepo, version)
	a.buildProvider(sourceCodeDir, providerName, version, customBuildCommand)

	name := extractRepoNameFromURL(gitRepo)
	a.moveBinaryToCorrectLocation(providerName, version, name)

	return true
}
