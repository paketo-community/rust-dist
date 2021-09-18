/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rust

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/heroku/color"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Cargo struct {
	Logger bard.Logger
}

func (c Cargo) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	c.Logger.Headerf("%s: %s to layer", color.BlueString(c.Name()), color.YellowString("Contributing"))

	AppendToPath(filepath.Join(layer.Path, "bin"))

	if err := os.Setenv("CARGO_HOME", layer.Path); err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to set $CARGO_HOME\n%w", err)
	}

	layer.BuildEnvironment.Override("CARGO_HOME", layer.Path)
	layer.LayerTypes = libcnb.LayerTypes{
		Build: true,
		Cache: true,
	}

	return layer, nil
}

func (c Cargo) Name() string {
	return "Cargo"
}
