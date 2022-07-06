package app

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"golang.org/x/mod/sumdb/dirhash"
)

const terraformLockfileName = ".terraform.lock.hcl"

type Lockfile struct {
	Provider []ProviderConfig `hcl:"provider,block"`
}

type ProviderConfig struct {
	Name        string   `hcl:"name,label"`
	Version     string   `hcl:"version"`
	Constraints string   `hcl:"constraints"`
	Hashes      []string `hcl:"hashes"`
}

func (a *App) UpgradeLockfile(inputLockfilePath string, outputLockfilePath string) {
	/*
		1. find lockfile
			parse the correct location
			man kann auch per parameter den genauen pfad an den command mitgeben
		2. parse lockfile
		3. find the entries that have local providers
			1. what to do when versions differ?
		4. calculate hash sum of local providers and replace the has of the entries
		5. Write HCL file back to original destination
	*/

	verifiedLockfilePath := getLockfile(inputLockfilePath)
	verifiedOutputPath := parseOutputLockfilePath(outputLockfilePath)
	log.Printf("Lockfile path: %s", verifiedLockfilePath)
	log.Printf("Output path: %s", verifiedOutputPath)

	// parse lockfile
	var config Lockfile
	err := hclsimple.DecodeFile(verifiedLockfilePath, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	for i, v := range config.Provider {
		filePath := a.Config.TerraformPluginDir + "/" + v.Name + "/" + v.Version + "/darwin_arm64"
		fmt.Println("filePath" + filePath)

		// if filePath is existent that means we have an local installed provider
		if isDirExistent(filePath) {
			hash := getCalculatedHashForProvider(filePath)
			// append the new hash to provide backwards compatability
			config.Provider[i].Hashes = append(v.Hashes, hash)
		}
	}
	newContents := createHclBody(config)
	writeFile(newContents, verifiedOutputPath)
}

// check if lockFilePath exists
// if not get lockfile from pwd directory
func getLockfile(lockFilePath string) string {
	if isDirExistent(lockFilePath) {
		return lockFilePath
	}
	return terraformLockfileName
}

func parseOutputLockfilePath(outputLockfilePath string) string {
	if outputLockfilePath == "" {
		return terraformLockfileName
	}
	return outputLockfilePath
}

func writeFile(contents string, path string) {
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(contents)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func createHclBody(hcl Lockfile) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	for _, v := range hcl.Provider {
		barBlock := rootBody.AppendNewBlock("provider", []string{v.Name})
		barBody := barBlock.Body()
		barBody.SetAttributeValue("version", cty.StringVal(v.Version))
		barBody.SetAttributeValue("constraints", cty.StringVal(v.Constraints))

		var listOfHashes []cty.Value
		for _, hash := range v.Hashes {
			listOfHashes = append(listOfHashes, cty.StringVal(hash))
		}
		list := cty.ListVal(listOfHashes)
		barBody.SetAttributeValue("hashes", list)
	}

	strToConvert := bytes.NewBuffer(f.Bytes()).String()
	return strToConvert
}

func getCalculatedHashForProvider(fullPath string) string {
	packageDir, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		log.Print("error")
	}
	// log.Print("dir" + packageDir)
	// The dirhash.HashDir result is already in our expected h1:...
	// format, so we can just convert directly to Hash.
	s, err := dirhash.HashDir(packageDir, "", dirhash.Hash1)
	if err != nil {
		log.Fatalf("Failed to load calculate hash: %s", err)
	}
	// log.Print("hash:" + s)
	return s
}
