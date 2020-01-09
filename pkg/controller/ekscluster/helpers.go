package ekscluster

import (
	"fmt"

	awsv1alpha1 "github.com/appvia/eks-operator/pkg/apis/aws/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	eks "github.com/aws/aws-sdk-go/service/eks"
)

// VerifyCredentials is responsible for verifying AWS creds
func VerifyCredentials(credentials *awsv1alpha1.AWSCredential) error {
	return nil
}

// Get an AWS session
func GetAWSSession(cred *awsv1alpha1.AWSCredential, region string) (*session.Session, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(cred.Spec.AccessKeyId, cred.Spec.SecretAccessKey, ""),
	})
	return sesh, err
}

// Get an EKS service
func GetEKSService(sesh *session.Session) (svc *eks.EKS, err error) {
	svc = eks.New(sesh)
	return svc, err
}

// Create an EKS cluster
func CreateEKSCluster(svc *eks.EKS, input *eks.CreateClusterInput) (output *eks.CreateClusterOutput, err error) {
	output, err = svc.CreateCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				fmt.Println(eks.ErrCodeResourceLimitExceededException, aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				fmt.Println(eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				fmt.Println(eks.ErrCodeUnsupportedAvailabilityZoneException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	return
}

// Describe a cluster
func DescribeEKSCluster(svc *eks.EKS, input *eks.DescribeClusterInput) (output *eks.DescribeClusterOutput, err error) {
	output, err = svc.DescribeCluster(input)
	return output, err
}

// Describe a cluster
func GetEKSClusterStatus(svc *eks.EKS, input *eks.DescribeClusterInput) (status string, err error) {
	cluster, err := svc.DescribeCluster(input)
	return *cluster.Cluster.Status, err
}

//	Lists all EKS clusters
func ListEKSClusters(svc *eks.EKS, input *eks.ListClustersInput) (output *eks.ListClustersOutput, err error) {
	output, err = svc.ListClusters(input)
	return output, err
}

// Check that a cluster exists
func EKSClusterExists(svc *eks.EKS, clusterName string) (exists bool, err error) {
	clusterList, err := svc.ListClusters(&eks.ListClustersInput{})
	if err != nil {
		return false, err
	}
	for _, i := range clusterList.Clusters {
		if clusterName == *i {
			return true, nil
		}
	}
	return false, nil
}

// Delete an EKS cluster
func DeleteEKSCluster(svc *eks.EKS, input *eks.DeleteClusterInput) (output *eks.DeleteClusterOutput, err error) {
	output, err = svc.DeleteCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				fmt.Println(eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceNotFoundException:
				fmt.Println(eks.ErrCodeResourceNotFoundException, aerr.Error())
			case eks.ErrCodeClientException:
				fmt.Println(eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				fmt.Println(eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				fmt.Println(eks.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	return
}
