package main

import (
    "fmt"
    "encoding/json"
    "path"
    "io/ioutil"
    "os"

    "../repo"
)

func CreateRepo(repoDir string) CommandError {
    //fileMode := (os.FileMode) 644
    fileMode := os.ModePerm

    // Try to create the repository directory. If it already exists, this should cause an error. We shouldn't try to create a repository in a directory that already exists.
    if err := os.Mkdir(repoDir, fileMode); err != nil {
        if os.IsExist(err) {
            // Tell the user we can't overwrite an existing repository.
            return CausedError(fmt.Sprintf("Can't create repository at %s because the directory already exists. Cannot create a repository in an existing directory.", repoDir), 11, err)
        } else if os.IsNotExist(err) {
            // Tell the user that the repository's parent directory probably doesn't exist.
            return CausedError(fmt.Sprintf("Can't create repository at %s. Make sure the parent directory exists.", repoDir), 12, err)
        } else {
            // An unknown error occurred.
            return CausedError(fmt.Sprintf("Failed to create repository at %s. An unknown error occurred.", repoDir), 10, err)
        }
    }

    // Determine the path to the index file.
    indexFilePath := path.Join(repoDir, repo.IndexFileName)

    // Get a new, blank index data struct.
    indexData := repo.NewBlankIndex()

    // Serialize the index structure to JSON...
    jsonData, jsonError := json.Marshal(indexData)

    if jsonError != nil {
        return CausedError("Failed to marshal index data to JSON. This probably shouldn't happen...", -1, jsonError)
    }
    
    // ...and write it to the index file.
    writeError := ioutil.WriteFile(indexFilePath, jsonData, fileMode)

    if writeError != nil {
        return CausedError("Failed to write index file.", 20, writeError)
    }

    return nil
}

