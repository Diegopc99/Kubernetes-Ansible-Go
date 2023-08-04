package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/Kubernetes-Ansible-Go/AnsibleAPI"
	"github.com/Kubernetes-Ansible-Go/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type ConfigDeploymentData struct {
	Namespace      string              `json:"namespace"`
	Name           string              `json:"name" binding:"required"`
	Replicas       int32               `json:"replicas"`
	ContainerName  string              `json:"containerName" binding:"required"`
	Image          string              `json:"image" binding:"required"`
	ContainerPorts []ContainerPortData `json:"containerPorts"`
	Labels         map[string]string   `json:"labels"`
}

type ContainerPortData struct {
	Name     string `json:"name" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
	Port     int32  `json:"port" binding:"required"`
}

func isValidProtocol(protocol string) bool {
	// A valid protocol can be "TCP", "UDP", or "SCTP".
	validProtocolRegex := regexp.MustCompile(`^(TCP|UDP|SCTP)$`)
	return validProtocolRegex.MatchString(strings.ToUpper(protocol))
}

func isValidName(name string) bool {
	// A valid name should consist of lowercase alphanumeric characters or '-', and it should start and end with an alphanumeric character.
	validNameRegex := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	return validNameRegex.MatchString(name)
}

func CreateDeployment(c *gin.Context) {

	var configDeployment ConfigDeploymentData

	// Process request params with required field validation
	if err := c.BindJSON(&configDeployment); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.Error(&exceptions.MissingFields{})
			return
		}
		c.Error(&exceptions.InvalidRequest{})
		return
	}

	if isValidName(configDeployment.ContainerName) == false {
		c.Error(&exceptions.InvalidRequest{})
		return
	}

	for _, portData := range configDeployment.ContainerPorts {
		fmt.Println(portData.Port)
		if isValidProtocol(portData.Protocol) == false {
			c.Error(&exceptions.InvalidRequest{})
			return
		}
	}

	yamlData, err := yaml.Marshal(configDeployment)
	if err != nil {
		c.Error(&exceptions.InternalError{})
		return
	}

	playbookFilename := "temp.yaml"

	// Save the YAML data to a file
	err = ioutil.WriteFile(playbookFilename, yamlData, 0644)
	if err != nil {
		c.Error(&exceptions.InternalError{})
		return
	}

	AnsibleAPI.ExecuteAnsible(playbookFilename)

	c.JSON(http.StatusOK, gin.H{})

}
