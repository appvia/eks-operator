package ekscluster

import (
	"fmt"

	awsv1alpha1 "github.com/appvia/eks-operator/pkg/apis/aws/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	eks "github.com/aws/aws-sdk-go/service/eks"
)

// VerifyCredentials is responsible for verifying AWS creds
func VerifyCredentials(credentials *awsv1alpha1.AWSCredential) error {
	return nil
}

// Get an AWS session
func GetAWSSession(credentials *awsv1alpha1.AWSCredential) (*session.Session, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	return sesh, err
}

// Get an EKS service
func GetEKSService(credentials *awsv1alpha1.AWSCredential) (svc *eks.EKS, err error) {
	svc = eks.New(session.New())
	return svc, err
}

// Create an EKS cluster
func CreateCluster(svc *eks.EKS, input *eks.CreateClusterInput) (output *eks.CreateClusterOutput, err error) {
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

// Delete an EKS cluster
func DeleteCluster(svc *eks.EKS, input *eks.DeleteClusterInput) (output *eks.DeleteClusterOutput, err error) {
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
