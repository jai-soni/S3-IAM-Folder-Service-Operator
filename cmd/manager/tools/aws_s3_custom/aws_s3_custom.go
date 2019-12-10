package aws_s3_custom

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
)

// PolicyDocument is our definition of our policies to be uploaded to IAM.
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

// StatementEntry will dictate what this policy will allow or not allow.
type StatementEntry struct {
	Effect   string
	Action   []string
	Resource string
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// CreateFolderIfNotExist Does somthing
func CreateFolderIfNotExist(filename, bucket, region string) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	_, err = svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == s3.ErrCodeNoSuchKey {
		result, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})

		if err != nil {
			fmt.Println("Create Folder error", err)
			return
		}

		fmt.Println("Success", result)
	} else {
		fmt.Println("GetFolder Error", err)
	}
}

// CreateUserIfNotExist Does somthing
func CreateUserIfNotExist(userName, region string) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})

	// Create a IAM service client.
	svc := iam.New(sess)

	_, err = svc.GetUser(&iam.GetUserInput{
		UserName: &os.Args[1],
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		result, err := svc.CreateUser(&iam.CreateUserInput{
			UserName: &userName,
		})

		if err != nil {
			fmt.Println("CreateUser Error", err)
			return
		}

		fmt.Println("Success", result)
	} else {
		fmt.Println("GetUser Error", err)
	}
}

// CreatePolicyIfNotExist does something
func CreatePolicyIfNotExist(filename, bucket, region, userName string) {

	var policyName = filename[:len(filename)-2] + "_s3_policy"
	var arnString = "arn:aws:s3:::" + bucket + "/" + filename[:len(filename)-2]
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})

	// Create a IAM service client.
	svc := iam.New(sess)

	// Builds our policy document for IAM.
	policy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			StatementEntry{
				Effect: "Allow",
				// Allows for DeleteItem, GetItem, PutItem, Scan, and UpdateItem
				Action: []string{
					"s3:*",
				},
				Resource: arnString,
			},
		},
	}

	b, err := json.Marshal(&policy)
	if err != nil {
		fmt.Println("Error marshaling policy", err)
		return
	}

	result, err := svc.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(b)),
		PolicyName:     aws.String(policyName),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("New policy", result)
	_, err = svc.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: result.Policy.Arn,
		UserName:  &userName,
	})

	if err != nil {
		fmt.Println("Unable to attach role policy to user", err)
		return
	}
	fmt.Println("Policy attached to user successfully")

}
