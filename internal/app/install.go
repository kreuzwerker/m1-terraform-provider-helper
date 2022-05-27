package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const requestTimeoutSeconds int = 2
const three int = 3
const four int = 4

type Provider struct {
	Repo        string `json:"source"`
	Description string `json:"description"`
}

type BuildCommandInformation struct {
	command         string
	startingVersion int
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
		log.Fatalf("Bash code did not run successfully: %s", err)
	}
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

	if len(version) > 0 {
		log.Println("version: " + version)
		ref, _ := r.ResolveRevision(plumbing.Revision(version))
		err = w.Checkout(&git.CheckoutOptions{
			Hash: *ref,
		})
		CheckIfError(err)
	} else {
		log.Println("No version specified, pulling and checking out main branch")
		executeBashCommand("git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@' | xargs git checkout && git pull", fullPath)
	}

	return repoDir
}

func extractRepoNameFromURL(url string) string {
	parts := strings.Split(url, "/")

	return parts[len(parts)-1]
}

func extractMajorVersionAsNumber(version string) int {
	sampleRegexp := regexp.MustCompile(`\d`)

	result := sampleRegexp.FindString(version)
	number, _ := strconv.Atoi(result)

	return number
}

func normalizeSemver(version string) string {
	if strings.HasPrefix(version, "v") {
		return version[1:]
	}

	return version
}

func createBuildCommand(providerName string, version string, goPath string) string {
	majorVersionNumberAsInt := extractMajorVersionAsNumber(version)

	buildCommands := make(map[string][]BuildCommandInformation)
	buildCommands["default"] = []BuildCommandInformation{{command: "make build", startingVersion: 0}}
	buildCommands["hashicorp/helm"] = []BuildCommandInformation{{command: "make build && cp terraform-provider-helm " + goPath + "/bin/" + "terraform-provider-helm", startingVersion: 0}}
	buildCommands["hashicorp/google"] = []BuildCommandInformation{{command: "gofmt -s -w ./tools.go  && make build", startingVersion: 0}}
	buildCommands["hashicorp/aws"] = []BuildCommandInformation{
		{command: "go get -u golang.org/x/sys && make tools && make fmt && gofmt -s -w ./tools.go && make build", startingVersion: 0},
		{command: "cd tools && go get -d github.com/pavius/impi/cmd/impi && cd .. && make tools && make build", startingVersion: three},
		{command: "make tools && make build", startingVersion: four},
	}

	buildCommandMap, exists := buildCommands[providerName]

	if exists {
		var foundBuildCommand string

		for _, v := range buildCommandMap {
			if majorVersionNumberAsInt >= v.startingVersion {
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

	filePath := a.Config.TerraformPluginDir + "/registry.terraform.io/" + providerName + "/" + version + "/darwin_arm64"
	err := os.MkdirAll(filePath, 0777)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "GOPATH: %s\n", a.Config.GoPath)
	pathOfExecutable := a.Config.GoPath + "/bin/" + executableName
	newPath := filePath + "/" + executableName + "_" + version + "_x5"

	log.Print("Move from " + pathOfExecutable + " to " + newPath)
	err = os.Rename(pathOfExecutable, newPath)

	if err != nil {
		log.Fatal(err)
	}
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
