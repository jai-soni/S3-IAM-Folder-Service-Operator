package controller

import (
	"github.com/sreeragsreenath/team2-kubeop/pkg/controller/folderservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, folderservice.Add)
}
