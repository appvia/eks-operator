package controller

import (
	"github.com/appvia/eks-operator/pkg/controller/ekscluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ekscluster.Add)
}
