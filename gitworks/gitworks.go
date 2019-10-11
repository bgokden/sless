package gitworks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func CloneBuildLoad(url, folder, username, password string) error {

	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("./temp", "clone")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Dir: %v\n", dir)

	defer os.RemoveAll(dir) // clean up

	_, gitErr := git.PlainClone(dir, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth: &http.BasicAuth{
			Username: username, // yes, this can be anything except an empty string
			Password: password,
		},
		URL:      url,
		Progress: os.Stdout,
	})

	if gitErr != nil {
		return gitErr
	}

	fullFolderPath := fmt.Sprint(dir + "/" + folder + "/*")
	objectPath := fmt.Sprint(dir + "/" + folder + "/" + "library.o")

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", objectPath, fullFolderPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmdErr := cmd.Run()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
	fmt.Printf("in all caps: %q\n", out.String())

	return nil
}
