// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	BaseBinaryVersion = "v0.8"
	BinaryName        = "unifi2mqtt"
	BaseDockerAccount = "mannkind"
)

// Everything below here is the same between projects
var BinaryArchs = []string{"amd64", "arm32v6", "arm64v8"}
var BinaryVersion = BaseBinaryVersion + time.Now().Format(".06002.1504")
var DockerImage = fmt.Sprintf("%s/%s", BaseDockerAccount, BinaryName)

// Commands
var g0 = sh.RunCmd("go")
var git = sh.RunCmd("git")
var docker = sh.RunCmd("docker")

type Go mg.Namespace
type Git mg.Namespace
type Docker mg.Namespace

// Default target to run when none is specified
var Default = All

// go:wire, go:format, go:vet, go:build, and go:test in that order
func All() {
	mg.SerialDeps(Go.Wire)
	mg.SerialDeps(Go.Format)
	mg.SerialDeps(Go.Vet)
	mg.SerialDeps(Go.Build)
	mg.SerialDeps(Go.Test)
	mg.SerialDeps(Go.Tidy)
}

// Remove the binary and architecture specific Dockerfiles
func Clean() error {
	fmt.Println("Cleaning")
	if err := sh.Run("rm", "-f", BinaryName); err != nil {
		return err
	}

	for _, arch := range BinaryArchs {
		if err := sh.Run("rm", "-f", "Dockerfile."+arch); err != nil {
			return err
		}
	}

	return nil
}

// Compile the application with the proper ldflags
func (Go) Build() error {
	fmt.Println("Building")
	return g0("build", "-ldflags", "-X \"main.Version="+BinaryVersion+"\" -X \"main.Name="+BinaryName+"\"", "-o", BinaryName, ".")
}

// Ensure the code is formatted properly
func (Go) Format() error {
	fmt.Println("Formatting")
	return g0("fmt", ".")
}

// Ensure the code passes vetting
func (Go) Vet() error {
	fmt.Println("Vetting")
	return g0("vet", ".")
}

// Get the compile-time DI tool
func (Go) GetWire() error {
	fmt.Println("Getting Wire")
	return g0("get", "github.com/google/wire/cmd/wire")
}

// Generate the dependencies at compile-time
func (Go) Wire() error {
	mg.SerialDeps(Go.GetWire)

	fmt.Println("Wiring")
	return sh.Run("wire", "gen")
}

// Run the application tests
func (Go) Test() error {
	mg.SerialDeps(Go.Build)

	fmt.Println("Testing")
	return g0("test", "--coverprofile", "/tmp/app.cover", "-v", ".")
}

// Run the tidy command
func (Go) Tidy() error {
	mg.SerialDeps(Go.Build)

	fmt.Println("Tidying")
	return g0("mod", "tidy")
}

// Create a new tag for the git repository using the generated binary version
func (Git) Tag() error {
	fmt.Println("Tagging Git Repo")
	return git("tag", "-f", BinaryVersion)
}

// Create a realese in Github by pushing the tags in the git repository
func (Git) Push() error {
	mg.SerialDeps(Git.Tag)

	fmt.Println("Pushing Git Repo Tags")
	if err := sh.Run("sed", "-i", "-e", "s/github.com/mannkind:$GITHUB_TOKEN@github.com/g", ".git/config"); err != nil {
		return err
	}

	if err := git("push", "--tags"); err != nil {
		return err
	}

	return nil
}

// Docker Build the docker images for each supported architecture
func (Docker) Build() error {
	for _, arch := range BinaryArchs {
		dockerfileWithArch := fmt.Sprintf("Dockerfile.%s", arch)
		dockerImageWithArch := fmt.Sprintf("%s:%s", DockerImage, arch)
		dockerImageWithArchVersion := fmt.Sprintf("%s-%s", dockerImageWithArch, BinaryVersion)
		dockerImageWithArchLatest := fmt.Sprintf("%s-latest", dockerImageWithArch)

		fmt.Println(fmt.Sprintf("Building image %s", dockerImageWithArchVersion))

		var golangArch = ""
		switch arch {
		case "amd64":
			golangArch = "amd64"
		case "arm32v6":
			golangArch = "arm"
		case "arm64v8":
			golangArch = "arm64"
		}

		if err := sh.Copy(dockerfileWithArch, "Dockerfile.template"); err != nil {
			return err
		}

		if err := sh.Run("sed", "-i", "-e", "s|__BASEIMAGE_ARCH__|"+arch+"|g", dockerfileWithArch); err != nil {
			return err
		}

		if err := sh.Run("sed", "-i", "-e", "s|__GOLANG_ARCH__|"+golangArch+"|g", dockerfileWithArch); err != nil {
			return err
		}

		if err := sh.Run("sed", "-i", "-e", "s|__BINARY_NAME__|"+BinaryName+"|g", dockerfileWithArch); err != nil {
			return err
		}

		if err := docker("build", "--pull", "-f", dockerfileWithArch, "-t", dockerImageWithArchVersion, "."); err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Tagging image %s as %s", dockerImageWithArchVersion, dockerImageWithArchLatest))
		if err := docker("tag", dockerImageWithArchVersion, dockerImageWithArchLatest); err != nil {
			return err
		}
	}

	mg.SerialDeps(Clean)

	return nil
}

// Upload the docker image to Dockerhub for each supported architecture
func (Docker) Push() error {
	mg.SerialDeps(Docker.Build)

	tags := []string{BinaryVersion, "latest"}
	for _, tag := range tags {
		// Must push the images before they can be used in a manifest
		for _, arch := range BinaryArchs {
			dockerImageWithTag := fmt.Sprintf("%s:%s-%s", DockerImage, arch, tag)
			fmt.Println(fmt.Sprintf("Pushing image %s ", dockerImageWithTag))
			if err := docker("push", dockerImageWithTag); err != nil {
				return err
			}
		}

		dockerManifest := fmt.Sprintf("%s:%s", DockerImage, tag)

		// Create the manifest
		fmt.Println(fmt.Sprintf("Creating multi-arch manifest %s", dockerManifest))
		if err := docker("manifest", "create", dockerManifest, fmt.Sprintf("%s:%s-%s", DockerImage, BinaryArchs[0], tag), fmt.Sprintf("%s:%s-%s", DockerImage, BinaryArchs[1], tag), fmt.Sprintf("%s:%s-%s", DockerImage, BinaryArchs[2], tag)); err != nil {
			return err
		}

		// Annotate the manifest
		for _, arch := range BinaryArchs {
			fmt.Println(fmt.Sprintf("Annotating multi-arch manifest %s-%s", dockerManifest, arch))

			var archName = ""
			var archVariant = ""
			switch arch {
			case "amd64":
				archName = ""
				archVariant = ""
			case "arm32v6":
				archName = "arm"
				archVariant = "v6"
			case "arm64v8":
				archName = "arm64"
				archVariant = "v8"
			}

			if archName == "" && archVariant == "" {
				continue
			}

			if err := docker("manifest", "annotate", dockerManifest, fmt.Sprintf("%s:%s-%s", DockerImage, arch, tag), "--os", "linux", "--arch", archName, "--variant", archVariant); err != nil {
				return err
			}
		}

		// Publish the manifest
		fmt.Println(fmt.Sprintf("Pushing multi-arch manifest %s", dockerManifest))
		if err := docker("manifest", "push", "--purge", dockerManifest); err != nil {
			return err
		}
	}
	return nil
}
