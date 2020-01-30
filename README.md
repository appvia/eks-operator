# AWS EKS Operator

## Overview
This operator can be used to create EKS clusters and node groups. This operator uses four CRDs:
* `AWSCredential` - defines an AWS access key/secret
* `AWSToken` - defines an STS token for authentication
* `EKSCluster` - defines an EKS cluster
* `EKSNodeGroup` - defines an EKS node group

## Prereqs
* Create an [EKS service IAM role](https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html#create-service-role) to assign to created clusters
* Create an [EKS worker IAM role](https://docs.aws.amazon.com/eks/latest/userguide/worker_node_IAM_role.html#create-worker-node-role)
* Activate STS in the regions you wish to create clusters [activate region](https://aws.amazon.com/blogs/security/aws-security-token-service-is-now-available-in-every-aws-region/)

## Use the operator
* Create an `EKSCredential` custom resource
* Create an `EKSCluster` referencing a credential
* Create an `EKSNodeGroup` referencing the EKS cluster

