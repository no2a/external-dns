/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Note: currently only supports IP targets (A records), not hostname targets
*/

package source

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/external-dns/endpoint"
)

// fakeSource is an implementation of Source that provides dummy endpoints for
// testing/dry-running of dns providers without needing an attached Kubernetes cluster.
type dirSource struct {
	dir string
}

// NewFakeSource creates a new fakeSource with the given config.
func NewDirSource(dir string) (Source, error) {
	return &dirSource{
		dir: "sourcedir",
	}, nil
}

func (sc *dirSource) AddEventHandler(ctx context.Context, handler func()) {
}

// Endpoints returns endpoint objects.
func (sc *dirSource) Endpoints(ctx context.Context) ([]*endpoint.Endpoint, error) {
	endpoints := make([]*endpoint.Endpoint, 0)
	files, err := os.ReadDir(sc.dir)
	if err != nil {
		fmt.Println("ReadDir err")
		panic(err)
	}
	for _, file := range files {
		var fp *os.File
		path := filepath.Join(sc.dir, file.Name())
		if fp, err = os.Open(path); err != nil {
			fmt.Println("Open err")
			panic(err)
		}
		defer fp.Close()
		reader := bufio.NewReader(fp)
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("ReadLine err")
			panic(err)
		}
		ep := endpoint.NewEndpoint(file.Name(), endpoint.RecordTypeA, string(line))
		endpoints = append(endpoints, ep)
	}
	return endpoints, nil
}
