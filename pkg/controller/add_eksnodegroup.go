package controller

import (
	"github.com/appvia/eks-operator/pkg/controller/eksnodegroup"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, eksnodegroup.Add)
}
