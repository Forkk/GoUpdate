// Copyright 2013 MultiMC Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type FileMD5Data struct {
	Path string
	MD5  string
}

func recursiveMD5Impl(root, currentPath string, current os.FileInfo, skipFiles []string) (data []FileMD5Data, err error) {
	err = nil
	data = []FileMD5Data{}

	relative, relErr := filepath.Rel(root, currentPath)
	if relErr != nil {
		err = relErr
		return
	}

	for _, skip := range skipFiles {
		if relative == skip {
			return
		}
	}

	//fmt.Printf("Path: %-40.40s Root: %-40.40s Relative to root: %-40.40s\n", currentPath, root, relative)
	if current.Mode().IsDir() {
		// If the file is a directory, get a list of its entries.
		if entries, readErr := ioutil.ReadDir(currentPath); readErr == nil {
			for _, entry := range entries {
				// Recurse! RECURSE! RECURSE!!!
				recurseData, recurseErr := recursiveMD5Impl(root, path.Join(currentPath, entry.Name()), entry, skipFiles)
				if recurseErr != nil {
					err = recurseErr
					return
				} else {
					data = append(data, recurseData...)
				}
			}
		} else {
			err = readErr
			return
		}
	} else if current.Mode().IsRegular() {
		file, openErr := os.Open(currentPath)

		if openErr != nil {
			err = openErr
			return
		}

		digest := md5.New()
		io.Copy(digest, file)
		md5Data := FileMD5Data{Path: relative, MD5: fmt.Sprintf("%x", digest.Sum(nil))}

		data = append(data, md5Data)
	}

	return
}

// Recursively calculates MD5 sums for all of the files in the given directory. Skips any files whose path relative to the directory matches a path in the skipFiles slice.
func RecursiveMD5Calc(path string, skipFiles []string) (data []FileMD5Data, err error) {
	if fileInfo, statErr := os.Stat(path); statErr == nil {
		data, err = recursiveMD5Impl(path, path, fileInfo, skipFiles)
		return
	} else {
		err = statErr
		data = nil
		return
	}
}
