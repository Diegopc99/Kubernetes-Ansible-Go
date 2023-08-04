package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Kubernetes-Ansible-Go/exceptions"
	"github.com/Kubernetes-Ansible-Go/handlers"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	/* appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry" *///
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// Structure matching the configuration file structure
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	Mode    string        `yaml:"mode"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func main() {

	// Define command-line flags
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Check if the config flag is provided
	if *configPath == "" {
		*configPath = "config/config.yaml"
	}

	// Read the contents of the YAML file
	yamlFile, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal("File " + *configPath + " not found")
	}

	// Parse the YAML file into the config struct
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting server ...")

	gin.SetMode(config.Mode)

	router := gin.Default()

	// Add custom error handlers
	router.NoRoute((func(c *gin.Context) {
		c.Error(&exceptions.EndpointNotFound{})
	}))
	router.NoMethod((func(c *gin.Context) {
		c.Error(&exceptions.MethodNotAllowed{})
	}))
	// Register the error handling middleware
	router.Use(exceptions.ErrorHandlerMiddleware())

	// Add endpoints
	router.POST("/deployment/create", handlers.CreateDeployment)
	//router.GET("/deployment/get", handlers.GetDeployment)
	//router.DELETE("/deployment/:deploymentName", handlers.DeleteDeployment)

	// Run server
	err = router.Run(config.Server.Host + ":" + config.Server.Port)
	if err != nil {
		panic(err)
	}

}
