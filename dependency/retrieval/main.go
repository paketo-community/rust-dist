// Copyright 2018-2026 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/paketo-buildpacks/libdependency/github"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/upstream"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

const (
	id       = "rust"
	name     = "Rust"
	purlName = "rust"

	org  = "rust-lang"
	repo = "rust"
)

func main() {
	retrieve.NewMetadataWithPlatforms(id, github.GetAllVersions(os.Getenv("GITHUB_TOKEN"), org, repo), generateMetadata)
}

func generateMetadata(versionFetcher versionology.VersionFetcher, platform retrieve.Platform) ([]versionology.Dependency, error) {
	rustArch := rustArch(platform.Arch)
	if rustArch == "" {
		fmt.Printf("Skipping %s/%s: unsupported architecture\n", platform.OS, platform.Arch)
		return nil, nil
	}

	versionString := versionFetcher.Version().String()

	source := fmt.Sprintf("https://static.rust-lang.org/dist/rustc-%s-src.tar.gz", versionString)
	sourceChecksum, err := upstream.GetSHA256OfRemoteFile(source)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", source, err)
	}

	uri := fmt.Sprintf("https://static.rust-lang.org/dist/rust-%s-%s-unknown-linux-gnu.tar.gz", versionString, rustArch)
	checksum, err := upstream.GetSHA256OfRemoteFile(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to checksum %s\n%w", uri, err)
	}

	dependency := cargo.ConfigMetadataDependency{
		Arch:     platform.Arch,
		Checksum: fmt.Sprintf("sha256:%s", checksum),
		CPE:      fmt.Sprintf("cpe:2.3:a:rust:rust:%s:*:*:*:*:*:*:*", versionString),
		ID:       id,
		Licenses: []interface{}{
			map[string]string{
				"type": "Apache-2.0",
				"uri":  "https://github.com/rust-lang/rust/blob/master/LICENSE-APACHE",
			},
			map[string]string{
				"type": "MIT",
				"uri":  "https://github.com/rust-lang/rust/blob/master/LICENSE-MIT",
			},
		},
		Name:           name,
		OS:             platform.OS,
		PURL:           purlWithArch(retrieve.GeneratePURL(purlName, versionString, checksum, uri), platform.Arch),
		SHA256:         checksum,
		Source:         source,
		SourceChecksum: fmt.Sprintf("sha256:%s", sourceChecksum),
		Stacks:         []string{"*"},
		URI:            uri,
		Version:        versionString,
	}

	d, err := versionology.NewDependency(dependency, fmt.Sprintf("%s-%s", platform.OS, platform.Arch))
	if err != nil {
		return nil, err
	}

	return []versionology.Dependency{d}, nil
}

// purlWithArch adds the arch qualifier to a PURL.
// libdependency's GeneratePURL does not set it.
func purlWithArch(purl string, arch string) string {
	if base, qualifiers, found := strings.Cut(purl, "?"); found {
		return fmt.Sprintf("%s?arch=%s&%s", base, arch, qualifiers)
	}
	return fmt.Sprintf("%s?arch=%s", purl, arch)
}

func rustArch(arch string) string {
	switch arch {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "aarch64"
	default:
		return ""
	}
}
