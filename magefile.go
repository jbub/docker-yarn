// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	dockerImage      = "jbub/docker-yarn"
	dockerBaseImage  = "node:9.7.1-alpine"
	dockerMaintainer = "Juraj Bubniak <juraj.bubniak@gmail.com>"
)

type versionInfo struct {
	Name       string
	Image      string
	Version    string
	Maintainer string
}

func (v versionInfo) tag(latest bool) string {
	if latest {
		return fmt.Sprintf("%v:latest", dockerImage)
	}
	return fmt.Sprintf("%v:%v", dockerImage, v.Name)
}

var versions = []versionInfo{
	{Name: "1.0", Version: "1.0.2", Image: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "1.1", Version: "1.1.0", Image: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "1.2", Version: "1.2.1", Image: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "1.4", Version: "1.4.1", Image: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "1.5", Version: "1.5.1", Image: dockerBaseImage, Maintainer: dockerMaintainer},
	{Name: "1.6", Version: "1.6.0", Image: dockerBaseImage, Maintainer: dockerMaintainer},
}

var dockerfileTmplString = `
FROM {{ .Image }}
MAINTAINER {{ .Maintainer }}

ENV YARN_VERSION={{ .Version }}
ENV PATH /root/.yarn/bin:$PATH

RUN apk --no-cache add gnupg curl bash binutils tar \
  && touch /root/.bashrc \
  && curl -o- -L https://yarnpkg.com/install.sh | bash -s -- --version ${YARN_VERSION} \
  && apk del gnupg curl tar binutils

ENTRYPOINT [ "yarn" ]`

var dockerfileTmpl = template.Must(template.New("dockerfile").Parse(dockerfileTmplString))

var docker = sh.RunCmd("docker")

// Generate generates dockerfiles for all versions.
func Generate() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get wd: %v", err)
	}

	for _, info := range versions {
		dir := filepath.Join(wd, info.Name)
		if err := ensureDir(dir); err != nil {
			return err
		}

		if err := genDockerfile(dir, info); err != nil {
			return err
		}
	}

	latest := versions[len(versions)-1]
	if err := genDockerfile(wd, latest); err != nil {
		return err
	}

	return nil
}

// Docker builds and runs docker images for all versions.
func Docker() error {
	mg.Deps(Generate)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get wd: %v", err)
	}

	for _, info := range versions {
		dir := filepath.Join(wd, info.Name)
		if err := buildAndRunDocker(dir, info, false); err != nil {
			return fmt.Errorf("could not build/run docker: %v", err)
		}
	}

	latest := versions[len(versions)-1]
	if err := buildAndRunDocker(wd, latest, true); err != nil {
		return fmt.Errorf("could not build/run docker: %v", err)
	}

	return nil
}

// Push pushes built docker images do docker hub.
func Push() error {
	mg.Deps(Docker)

	for _, info := range versions {
		if err := pushDocker(info, false); err != nil {
			return fmt.Errorf("could not push docker image: %v", err)
		}
	}

	latest := versions[len(versions)-1]
	if err := pushDocker(latest, true); err != nil {
		return fmt.Errorf("could not push docker image: %v", err)
	}

	return nil
}

func pushDocker(info versionInfo, latest bool) error {
	return docker("push", info.tag(latest))
}

func buildAndRunDocker(dir string, info versionInfo, latest bool) error {
	if err := docker("build", "-t", info.tag(latest), dir); err != nil {
		return err
	}
	if err := docker("run", "--interactive", "--tty", "--rm", info.tag(latest)); err != nil {
		return err
	}
	return nil
}

func genDockerfile(dir string, info versionInfo) error {
	fpath := filepath.Join(dir, "Dockerfile")
	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer fp.Close()

	if err := dockerfileTmpl.Execute(fp, info); err != nil {
		return fmt.Errorf("could not execute template: %v", err)
	}
	return nil
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not stat path: %v", err)
		}

		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return fmt.Errorf("could not create dir: %v", err)
		}
	}
	return nil
}
