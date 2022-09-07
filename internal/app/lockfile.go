package app

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
	Constraints *string  `hcl:"constraints"`
	Hashes      []string `hcl:"hashes"`
}

/*
 1. find lockfile
    parse the correct location
    man kann auch per parameter den genauen pfad an den command mitgeben
 2. parse lockfile
 3. find the entries that have local providers
 1. what to do when versions differ?
 4. calculate hash sum of local providers and replace the has of the entries
 5. Write HCL file back to original destination.
*/
func (a *App) UpgradeLockfile(inputLockfilePath string, outputLockfilePath string) {
	verifiedLockfilePath := getLockfile(inputLockfilePath)
	verifiedOutputPath := parseOutputLockfilePath(outputLockfilePath)

	log.Printf("Lockfile path: %s", verifiedLockfilePath)
	log.Printf("Output path: %s", verifiedOutputPath)

	config := parseHclLockfile(verifiedLockfilePath)

	for i, v := range config.Provider {
		filePath := a.Config.TerraformPluginDir + "/" + v.Name + "/" + v.Version + "/darwin_arm64"

		// if filePath is existent that means we have an local installed provider
		if isDirExistent(filePath) {
			hash := getCalculatedHashForProvider(filePath)
			// append the new hash to provide backwards compatibility
			v.Hashes = append(v.Hashes, hash)
			config.Provider[i].Hashes = v.Hashes
		}
	}

	newContents := createHclBody(config)
	writeFile(newContents, verifiedOutputPath)
}

func parseHclLockfile(filepath string) Lockfile {
	var config Lockfile
	err := hclsimple.DecodeFile(filepath, nil, &config)

	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	return config
}

// check if lockFilePath exists
// if not get lockfile from pwd directory.
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
		log.Println(err)

		return
	}

	defer f.Close()

	_, err2 := f.WriteString(contents)

	if err2 != nil {
		log.Println(err)

		return
	}
}

func createHclBody(hcl Lockfile) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, v := range hcl.Provider {
		barBlock := rootBody.AppendNewBlock("provider", []string{v.Name})
		barBody := barBlock.Body()
		barBody.SetAttributeValue("version", cty.StringVal(v.Version))

		if v.Constraints != nil && *v.Constraints != "" {
			barBody.SetAttributeValue("constraints", cty.StringVal(*v.Constraints))
		}

		if len(v.Hashes) > 0 {
			hashToks := encodeHashSetTokens(v.Hashes)
			barBody.SetAttributeRaw("hashes", hashToks)
		}
	}

	convertedString := bytes.NewBuffer(f.Bytes()).String()
	fullString := "# This file is maintained automatically by \"terraform init\".\n# Manual edits may be lost in future updates.\n\n" + convertedString

	return fullString
}

// Taken from https://github.com/hashicorp/terraform/blob/aeefde7428b836646ba9622f1bb313e6dfe2ca87/internal/depsfile/locks_file.go#L454
// Based on this issue: https://github.com/hashicorp/hcl/issues/542
func encodeHashSetTokens(hashes []string) hclwrite.Tokens {
	// We'll generate the source code in a low-level way here (direct
	// token manipulation) because it's desirable to maintain exactly
	// the layout implemented here so that diffs against the locks
	// file are easy to read; we don't want potential future changes to
	// hclwrite to inadvertently introduce whitespace changes here.
	ret := hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenOBrack,
			Bytes: []byte{'['},
		},
		{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		},
	}

	// Although lock.hashes is a slice, we de-dupe and sort it on
	// initialization so it's normalized for interpretation as a logical
	// set, and so we can just trust it's already in a good order here.
	for _, hash := range hashes {
		hashVal := cty.StringVal(hash)
		ret = append(ret, hclwrite.TokensForValue(hashVal)...)
		ret = append(ret, hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenComma,
				Bytes: []byte{','},
			},
			{
				Type:  hclsyntax.TokenNewline,
				Bytes: []byte{'\n'},
			},
		}...)
	}
	ret = append(ret, &hclwrite.Token{
		Type:  hclsyntax.TokenCBrack,
		Bytes: []byte{']'},
	})

	return ret
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
