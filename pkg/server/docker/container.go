// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker_server

import (
	"fmt"

	v1 "github.com/nitrictech/boxygen/pkg/proto/builder/v1"
)

type ContainerState interface {
	Name() string
	Lines() []string
	// TODO: We will want to replace this with a more op/args model to translate better between more types of container file formats
	AddLine(line string)
	LogLine(line string) string
	AddDependency(name string)
	Dependencies() []string
	Ignore() []string
}

// ContainerState
type containerStateImpl struct {
	// unique name for this container state
	name string
	// container states that this container depends on
	dependsOn []string
	// lines composing this container image state
	lines []string
	// patterns to ignore when using file ops such as COPY
	ignore []string
}

func (c *containerStateImpl) Name() string {
	return c.name
}

func (c *containerStateImpl) Lines() []string {
	return c.lines
}

func (c *containerStateImpl) Ignore() []string {
	return c.ignore
}

func (c *containerStateImpl) AddLine(line string) {
	if c.lines == nil {
		c.lines = make([]string, 0)
	}

	c.lines = append(c.lines, line)
}

func (c *containerStateImpl) LogLine(line string) string {
	return fmt.Sprintf("Append [%s] to container %s", line, c.Name())
}

func (c *containerStateImpl) AddDependency(name string) {
	if c.dependsOn == nil {
		c.dependsOn = make([]string, 0)
	}

	c.dependsOn = append(c.dependsOn, name)
}

func (c *containerStateImpl) Dependencies() []string {
	if c.dependsOn == nil {
		c.dependsOn = make([]string, 0)
	}

	return c.dependsOn
}

func appendAndLog(line string, cs ContainerState, srv BuilderPbServer) {
	cs.AddLine(line)
	srv.Send(&v1.OutputResponse{
		Log: []string{cs.LogLine(line)},
	})
}
