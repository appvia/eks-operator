package main

import (
	"context"
	"fmt"

	awsv1alpha1 "github.com/appvia/eks-operator/pkg/apis/aws/v1alpha1"

	configv1 "github.com/appvia/hub-apis/pkg/apis/config/v1"
	"github.com/appvia/hub-apis/pkg/publish"
	hschema "github.com/appvia/hub-apis/pkg/schema"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	// awsClass is the provider class published into the hub
	awsClass = configv1.Class{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws",
			Namespace: "hub",
		},
		Spec: configv1.ClassSpec{
			APIVersion:    awsv1alpha1.SchemeGroupVersion.String(),
			AutoProvision: false,
			Category:      "cluster",
			Description:   "AWS EKS provides a means to provision Kubernetes clusters and node groups",
			DisplayName:   "AWS EKS",
			Requires: metav1.GroupVersionKind{
				Group:   awsv1alpha1.SchemeGroupVersion.Group,
				Kind:    "EKSCluster",
				Version: awsv1alpha1.SchemeGroupVersion.Version,
			},
			Plans: []string{},
			Resources: []configv1.ClassResource{
				{
					Group:            awsv1alpha1.SchemeGroupVersion.Group,
					Kind:             "EKSCluster",
					Plans:            []string{},
					DisplayName:      "EKS Cluster",
					ShortDescription: "Provisions an EKS Cluster",
					LongDescription:  "Provides the ability to provision an EKS cluster",
					Scope:            configv1.TeamScope,
					Version:          awsv1alpha1.SchemeGroupVersion.Version,
				},
				{
					Group:            awsv1alpha1.SchemeGroupVersion.Group,
					Kind:             "EKSNodeGroup",
					Plans:            []string{},
					DisplayName:      "EKS Node Group",
					ShortDescription: "Provisions an EKS Node Group",
					LongDescription:  "Provides the ability to provision an EKS node group",
					Scope:            configv1.TeamScope,
					Version:          awsv1alpha1.SchemeGroupVersion.Version,
				},
				{
					Group:            awsv1alpha1.SchemeGroupVersion.Group,
					Kind:             "AWSCredential",
					Plans:            []string{},
					DisplayName:      "AWS Credentials",
					ShortDescription: "The AWS credentials to use",
					LongDescription:  "The credentials used to provision cluters and node groups from",
					Scope:            configv1.TeamScope,
					Version:          awsv1alpha1.SchemeGroupVersion.Version,
				},
			},
			Schemas: schema.ConvertToJSON(),
		},
	}
)

// publishOperator is responsible for injecting the classes configuration
// into the hub-api and crds
func publishOperator(cfg *rest.Config) error {
	// @step: publish the CRDs in the hub
	ac, err := publish.NewExtentionsAPIClient(cfg)
	if err != nil {
		return err
	}

	if publishCRDs {
		if err := publish.ApplyCustomResourceDefinitions(ac, schema.GetCustomResourceDefinitions()); err != nil {
			return fmt.Errorf("failed to register the operator crds: %s", err)
		}
	}

	if publishClasses {
		c, err := hschema.NewClient(cfg)
		if err != nil {
			return err
		}

		return hschema.PublishAll(context.TODO(), c, awsClass, []configv1.Plan{})
	}

	return nil
}
