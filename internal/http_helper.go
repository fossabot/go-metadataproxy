package internal

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func remoteIP(addr string) string {
	return strings.Split(addr, ":")[0]
}

func findContainerRoleByAddress(addr string) (*iam.GetRoleOutput, error) {
	remoteIP := remoteIP(addr)

	container, err := findDockerContainer(remoteIP)
	if err != nil {
		return nil, err
	}

	roleName, err := findDockerContainerIAMRole(container)
	if err != nil {
		return nil, err
	}

	role, err := readRoleFromAWS(roleName)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func isCompatibleAPIVersion(r *http.Request) bool {
	vars := mux.Vars(r)
	return vars["api_version"] >= "2012-01-12"
}

func httpError(err error, w http.ResponseWriter, r *http.Request) {
	log.Error(err)
	http.NotFound(w, r)
}