package main

import (
    "bufio"
    "os"
    "path/filepath"
    log "github.com/Sirupsen/logrus"

    // Minimalist http framework
    "github.com/gin-gonic/gin"

    // Tar utilities for go
    "github.com/Rolinh/targo"

    // go client for git
    //"github.com/libgit2/git2go"

    // go client for docker
    "github.com/GoPex/dockerclient"
)

type Login struct {
    User     string `json:"user" binding:"required"`
}

func ping(c *gin.Context) {
    go build_from_directory("gopex/unleash_test_repository", "")

    c.JSON(200, gin.H{
        "message": "pong",
    })

}

func loginJSON(c *gin.Context) {
    var json Login
    if c.BindJSON(&json) == nil {
        log.Println(json.User)
    }
    log.Println(c)
}

    // Samalba docker client
    //docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

    // Listen to events
    //docker.StartMonitorEvents(eventCallback, nil)
//func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
    //log.Println("----------- NEW EVENT")
    //log.Println(event.Status)
    //log.Println(event.ID)
    //log.Println(event.Type)
    //log.Println(event.Action)
    //log.Println("----------- END OF EVENT")
//}

func build_from_directory(repository_name string, branch string) {
    log.Info("Building Dockerfile for repository ", repository_name)

    // Initialize a Docker client
    log.Info("woot")
    docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)
    log.Info("woot")

    // Path to the cloned repository to build
    repository_path := filepath.Join("/tmp", repository_name)
    log.Info("Path to the repository ", repository_path)

    // Create a tar from the repository
    targo.CreateInPlace(repository_path)
    repository_tar_path := repository_path + ".tar"
    log.Info("Path to the tar of repository ", repository_tar_path)

    // Send the tar to Docker for building
    dockerBuildContext, err := os.Open(repository_tar_path)
    defer dockerBuildContext.Close()

    // Construct the tag to use when building the image
    tag := ""
    if branch != ""{
        tag = ":" + branch
    } else {
        tag = ":latest"
    }
    log.Debug("Tag used for the image ", tag)

    // Build the image configuration send to Docker
    buildImageConfig := &dockerclient.BuildImage{
            Context:        dockerBuildContext,
            RepoName:       repository_name + tag,
            SuppressOutput: false,
    }

    // Send the build request to Docker
    reader, err := docker.BuildImage(buildImageConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Capture the build output and wait its end before continuing
    rd := bufio.NewScanner(reader)
    for rd.Scan() {
        log.Println(rd.Text())
    }

    // The Docker build image output has ended, check for errors
    err = rd.Err()
    if err != nil {
        log.Fatal("Read Error:", err)
    }

    log.Info("Repository %s has been builded !", repository_name)
}

func main() {

    // Create a default gin stack
    r := gin.Default()

    // Routes
    r.GET("/ping", ping)
    r.POST("/loginJSON", loginJSON)


    // Unleash!
    r.Run() // listen and serve on port defined by environment variable PORT
}


    //// To force the import
    //opts := git.CloneOptions{}
    //url := "https://github.com/GoPex/beekeeper_rails_app.git"
    //path := "/tmp/beekeeper_rails_app"
    //repo, err := git.Clone(url, path, &opts)
    //if err != nil {
        //log.Println(err)
    //} else {
        //log.Println(repo)
    //}

