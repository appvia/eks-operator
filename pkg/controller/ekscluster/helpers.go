package ekscluster

import (
	"fmt"
	"unicode/utf8"
	"regexp"

	strings "strings"
	b64 "encoding/base64"
	url "net/url"
	http "net/http"

	awsv1alpha1 "github.com/appvia/eks-operator/pkg/apis/aws/v1alpha1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"

	sts "github.com/aws/aws-sdk-go/service/sts"
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

func GetSTSService(sesh *session.Session) (svc *sts.STS, err error) {
	svc := sts.New(sesh)
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

// Check if a cluster exists
func CheckEKSClusterExists(svc *eks.EKS, input *eks.DescribeClusterInput) (exists bool, err error) {
	_, err = svc.DescribeCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				fmt.Println(aerr.Error())
				return false, err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return false, err
		}
	}
	return true, nil
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

func GetBearerToken(cred *awsv1alpha1.AWSCredential, clusterId, region string, expiration int64) (bearer string, error err) {
	signer = Signer{
		Credentials: credentials.NewStaticCredentials(cred.Spec.AccessKeyId, cred.Spec.SecretAccessKey, ""),
	}

	destUrl, err := url.Parse("https://sts." + region + ".amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15")

	if err != nil {
		return
	}

	header = map[string][]string{
		"x-k8s-aws-id": {clusterId},
	}

	request = http.Request{
		Method: "GET",
		URL: destUrl,
		Body: nil,
		Header: header,
	}

	// Presign the http request
	signedUrl, err := signer.Presign(request, nil, "sts", region, time.Duration(expiration*time.Second), time.Now())

	// Base64 encode it
	encodedUrl := b64.StdEncoding.EncodeToString([]byte(signedUrl))

	// Remove padding
	bearerKey := strings.ReplaceAll(encodedUrl, "=", "")

	bearer = "k8s-aws-v1." + bearerKey

	return bearer, err
}
