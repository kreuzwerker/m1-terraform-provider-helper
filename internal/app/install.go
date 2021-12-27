package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const requestTimeoutSeconds int = 2

type Provider struct {
	Repo        string `json:"source"`
	Description string `json:"description"`
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

func (a *App) getProviderData(providerName string) Provider {
	url := "https://registry.terraform.io/v1/providers/" + providerName

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * time.Duration(float64(requestTimeoutSeconds)))
		println("Cancel")
		cancel()
	}()

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		panic(err.Error())
	}
	// #nosec G107
	res, err := client.Do(req)

	if err != nil {
		panic(err.Error())
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data Provider

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	return data
}

func (a *App) cloneRepo(gitURL string) {
	if !a.isDirExistent(a.Config.ProvidersCacheDir) {
		err := os.Mkdir(a.Config.ProvidersCacheDir, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	command := "cd " + a.Config.ProvidersCacheDir + " && git clone " + gitURL
	log.Println("Executing" + command)
	bashCmd := exec.Command("sh", "-c", command)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	bashCmd.Stdout = mw
	bashCmd.Stderr = mw

	if err := bashCmd.Run(); err != nil {
		log.Fatalf("Bash code did not run successfully: %s", err)
	}

	log.Println(stdBuffer.String())
}

// if repo is not check out yet
//		- clone to cli dir
// if already exists: dont clone, simply cd
// on both casees: checkout version
// return path to dir.
func (a *App) checkoutSourceCode(gitURL string, version string) string {
	var r *git.Repository

	repoDir := strings.Split(gitURL, "/")[1]
	path := a.Config.ProvidersCacheDir + "/" + repoDir

	if !a.isDirExistent(path) {
		a.cloneRepo(gitURL)
	}

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// Clean the repository
	executeBashCommand("git reset --hard && git clean -d -f -q", path)

	if len(version) > 0 {
		log.Println("version: " + version)
		ref, _ := r.ResolveRevision(plumbing.Revision(version))
		err = w.Checkout(&git.CheckoutOptions{
			Hash: *ref,
		})
		CheckIfError(err)
	} else {
		log.Println("No version specified, pulling and checking out main branch")
		executeBashCommand("git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@' | xargs git checkout && git pull", path)
	}

	return repoDir
}

func (a *App) createBuildCommand(providerName string) string {
	buildCommands := make(map[string]string)
	buildCommands["default"] = "make build"
	buildCommands["hashicorp/aws"] = "cd tools && go get -d github.com/pavius/impi/cmd/impi && cd .. && make tools && make build"

	buildCommand, exists := buildCommands[providerName]

	if exists {
		return buildCommand
	}

	return buildCommands["default"]
}

func (a *App) buildProvider(dir string, providerName string) {
	buildCommand := a.createBuildCommand(providerName)
	// #nosec G204
	bashCmd := exec.Command("sh", "-c", "cd "+a.Config.ProvidersCacheDir+"/"+dir+" && "+buildCommand)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	bashCmd.Stdout = mw
	bashCmd.Stderr = mw

	if err := bashCmd.Run(); err != nil {
		log.Fatalf("Bash code did not run successfully: %s", err)
	}

	log.Println(stdBuffer.String())
}

func (a *App) moveBinaryToCorrectLocation(providerName string, version string, executableName string) {
	if len(version) == 0 {
		version = "master"
	}

	filePath := a.Config.TerraformPluginDir + "/registry.terraform.io/" + providerName + "/" + version + "/darwin_arm64"
	err := os.MkdirAll(filePath, 0777)

	if err != nil {
		log.Fatal(err)
	}

	pathOfExecutable := a.Config.GoPath + "/bin/" + executableName
	newPath := filePath + "/" + executableName + "_" + version + "_x5"

	log.Print("Move from " + pathOfExecutable + " to " + newPath)
	err = os.Rename(pathOfExecutable, newPath)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) Install(providerName string, version string) bool {
	providerData := a.getProviderData(providerName)
	fmt.Fprintf(os.Stdout, "Repo: %s\n", providerData.Repo)

	gitRepo := strings.Replace(providerData.Repo, "https://github.com/", "git@github.com:", 1)
	fmt.Fprintf(os.Stdout, "GitRepo: %s\n", gitRepo)

	sourceCodeDir := a.checkoutSourceCode(gitRepo, version)
	a.buildProvider(sourceCodeDir, providerName)

	name := strings.Split(gitRepo, "/")[1]
	a.moveBinaryToCorrectLocation(providerName, version, name)

	return true
}
