package serverworks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"plugin"

	git "gopkg.in/src-d/go-git.v4"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type ServerFunctionConf struct {
	Sourceurl string `yaml:"sourceurl,omitempty" json:"sourceurl,omitempty"`
	SubFolder string `yaml:"subfolder,omitempty" json:"subfolder,omitempty"`
	Path      string `yaml:"path,omitempty" json:"path,omitempty"`
	Username  string `yaml:"username,omitempty" json:"username,omitempty"`
	Password  string `yaml:"password,omitempty" json:"password,omitempty"`
}

type ServerWorksConf struct {
	Serverfunctions map[string]ServerFunctionConf `yaml:"serverfunctions,omitempty" json:"serverfunctions,omitempty"`
}

type ServerWorks struct {
	Conf *ServerWorksConf
	Mux  *http.ServeMux
}

var versionEnvVar string
var IsReady bool

func ready(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ready\n")
}

func healty(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "healty\n")
}

func index(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "OK\n")
}

func version(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, fmt.Sprintf("%v\n", versionEnvVar))
}

type Server interface {
	Name() string
	Serve(http.ResponseWriter, *http.Request)
}

func (s *ServerWorks) Load(path, mod string) {

	// determine module to load

	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	symServer, err := plug.Lookup("Server")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type Greeter (defined above)
	var server Server
	server, ok := symServer.(Server)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	fmt.Printf("Handle on path %v %v\n", path, server.Name())
	s.Mux.HandleFunc(path, server.Serve)
}

// const TickPeriod = 30 * time.Second

func (s *ServerWorks) loader() {
	for k, v := range s.Conf.Serverfunctions {
		log.Printf("Loading %v\n", k)
		s.CloneBuildLoad(v.Sourceurl, v.SubFolder, v.Username, v.Password, v.Path)
		log.Printf("Loaded %v\n", k)
	}
}

func (s *ServerWorks) RunServe() {
	IsReady = false

	s.Mux = http.NewServeMux()

	versionEnvVar = os.Getenv("VERSION")

	s.Mux.HandleFunc("/healty", healty)
	s.Mux.HandleFunc("/ready", ready)
	s.Mux.HandleFunc("/version", version)

	if s.Conf != nil {
		s.loader()
	}

	IsReady = true

	err := http.ListenAndServe(":8080", s.Mux)
	log.Fatal(err)
}

func (s *ServerWorks) CloneBuildLoad(url, folder, username, password, path string) error {

	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("./tempdata", "clone")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Dir: %v\n", dir)

	defer os.RemoveAll(dir) // clean up

	_, gitErr := git.PlainClone(dir, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth: &githttp.BasicAuth{
			Username: username, // yes, this can be anything except an empty string
			Password: password,
		},
		URL:      url,
		Progress: os.Stdout,
	})

	if gitErr != nil {
		fmt.Printf("Error: %v\n", gitErr)
		return gitErr
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullFolderPath := fmt.Sprint(currentDir + "/" + dir + "/" + folder + "/main.go")
	objectPath := fmt.Sprint(currentDir + "/" + dir + "/" + folder + "/" + "library.o")

	execpath, err := exec.LookPath("go")
	if err != nil {
		log.Fatalf("Path error %v\n", err)
	}
	fmt.Printf("%v\n", execpath)

	commandText := fmt.Sprintf("%v %v %v %v %v %v", execpath, "build", "-buildmode=plugin", "-o", objectPath, fullFolderPath)
	fmt.Printf("%v\n", commandText)
	cmd := exec.Command(execpath, "build", "-buildmode=plugin", "-o", objectPath, fullFolderPath)
	// fmt.Printf("os.Environ(): %v", os.Environ())
	cmd.Env = append(
		os.Environ(),
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmdErr := cmd.Run()

	fmt.Printf("out: %v err: %v\n", out.String(), stderr.String())
	if cmdErr != nil {
		fmt.Printf("Error: %v\n", cmdErr)

		log.Fatal(cmdErr)
	}

	s.Load(path, objectPath)

	return nil
}
