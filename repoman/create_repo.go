package main

import (
    "fmt"
    "encoding/json"
    "path"
    "io/ioutil"
    "os"
    "errors"

    "../repo"
)

func CreateRepo(repoDir string) error {
    //fileMode := (os.FileMode) 644
    fileMode := os.ModePerm

    // Determine the path to the index file.
    indexFilePath := path.Join(repoDir, repo.IndexFileName)

    // First, check and make sure that there isn't already a blank repository at the given path.
    if _, err := os.Stat(indexFilePath); !os.IsNotExist(err) {
        // If the file already exists, abort the operation with an error.
        return errors.New(fmt.Sprintf("There is already an existing GoUpdate repository at %s.", repoDir))
    }

    // If the repo directory doesn't exist, create it.
    if _, err := os.Stat(repoDir); os.IsNotExist(err) {
        // Try to create the directory. If this fails, return an error.
        createDirError := os.Mkdir(repoDir, fileMode)

        if createDirError != nil {
            return createDirError
        }
    }

    // Get a new, blank index data struct.
    indexData := repo.NewBlankIndex()

    // Serialize the index structure to JSON...
    jsonData, jsonError := json.Marshal(indexData)

    if jsonError != nil {
        return jsonError
    }
    
    // ...and write it to the index file.
    writeError := ioutil.WriteFile(indexFilePath, jsonData, fileMode)

    if writeError != nil {
        return writeError
    }

    return nil
}

