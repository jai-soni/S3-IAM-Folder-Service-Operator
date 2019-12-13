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
	"github.com/aws/aws-sdk-go/service/sts"
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
func CreateFolderIfNotExist(accessKeyID, secretAccessKey, filename, bucketName, region string) (success bool) {
	fmt.Println(">>>>>>>>>>>>>>>>")
	fmt.Println("Key %v", accessKeyID)
	fmt.Println("Secret %v", secretAccessKey)
	fmt.Println("BucketName %v", bucketName)
	fmt.Println("Region %v", region)
	fmt.Println(">>>>>>>>>>>>>>>>")
	success = false
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	_, err = svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == s3.ErrCodeNoSuchKey {
		result, err := svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(filename),
		})

		if err != nil {
			fmt.Println("Create Folder error", err)
			return
		}

		fmt.Println("Success", result)
		success = true
		return
	} else {
		fmt.Println("GetFolder Error", err)
		return
	}
}

// CreateUserIfNotExist Does somthing
func CreateUserIfNotExist(accessKeyID, secretAccessKey, userName, region string) (awsAccessKey string, awsSecretAccessKey string, success bool) {
	success = false
	awsAccessKey = ""
	awsSecretAccessKey = ""

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	// Create a IAM service client.
	svc := iam.New(sess)

	_, err = svc.GetUser(&iam.GetUserInput{
		UserName: &userName,
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
		result, err := svc.CreateUser(&iam.CreateUserInput{
			UserName: &userName,
		})

		if err != nil {
			fmt.Println("CreateUser Error", err)
			return
		}

		accessKeyResult, accessKeyErr := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
			UserName: aws.String(userName),
		})

		if accessKeyErr != nil {
			fmt.Println("Error", accessKeyErr)
			return
		}

		fmt.Println("Username created :", *result.User.UserName)
		success = true
		awsAccessKey = *accessKeyResult.AccessKey.AccessKeyId
		awsSecretAccessKey = *accessKeyResult.AccessKey.SecretAccessKey
		// fmt.Println("Secrets for User", *accessKeyResult.AccessKey)
		fmt.Println("awsAccessKey :", awsAccessKey)
		fmt.Println("awsSecretAccessKey :", awsSecretAccessKey)
		return
	}
	return
}

// CreateKeyIfNotExist is used to create key if not exists
func CreateKeyIfNotExist(accessKeyID, secretAccessKey, userName, region string) (awsAccessKey string, awsSecretAccessKey string, success bool) {
	success = false
	awsAccessKey = ""
	awsSecretAccessKey = ""

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	// Create a IAM service client.
	svc := iam.New(sess)
	result, err := svc.ListAccessKeys(&iam.ListAccessKeysInput{
		MaxItems: aws.Int64(5),
		UserName: aws.String(userName),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", result)

	for _, b := range result.AccessKeyMetadata {
		fmt.Printf("* %s access key deleted\n",
			aws.StringValue(b.AccessKeyId))
		svc.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: b.AccessKeyId,
			UserName:    &userName,
		})
	}
	accessKeyResult, accessKeyErr := svc.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(userName),
	})

	if accessKeyErr != nil {
		fmt.Println("Error", accessKeyErr)
		return
	}

	success = true
	awsAccessKey = *accessKeyResult.AccessKey.AccessKeyId
	awsSecretAccessKey = *accessKeyResult.AccessKey.SecretAccessKey
	// fmt.Println("Secrets for User", *accessKeyResult.AccessKey)
	fmt.Println("awsAccessKey :", awsAccessKey)
	fmt.Println("awsSecretAccessKey :", awsSecretAccessKey)
	return
}

func getUserAccountNumber(accessKeyID, secretAccessKey, region string) (accountNumber string) {
	accountNumber = ""
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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
	accountNumber = *result.Account
	return
}

func attachPolicyToUser(accessKeyID, secretAccessKey, region, policyArn, userName string) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	svc := iam.New(sess)

	_, err = svc.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: &policyArn,
		UserName:  &userName,
	})

	if err != nil {
		fmt.Println("Unable to attach role policy to user", err)
		return
	}
}

// CreatePolicyIfNotExist does something
func CreatePolicyIfNotExist(accessKeyID, secretAccessKey, filename, bucket, region, userName string) (success bool) {
	success = false
	var policyName = filename[:len(filename)-1] + "_s3_policy"
	var arnString = "arn:aws:s3:::" + bucket + "/" + filename[:len(filename)-1]
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})

	// Create a IAM service client.
	svc := iam.New(sess)

	// Check if the policy exists
	// userPolicyArn := "arn:aws:iam::aws:policy/" + policyName
	userPolicyArn := "arn:aws:iam::" + getUserAccountNumber(accessKeyID, secretAccessKey, region) + ":policy/" + policyName
	result2, err := svc.GetPolicy(&iam.GetPolicyInput{
		PolicyArn: &userPolicyArn,
	})

	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
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
		success = true
		fmt.Println("Policy created and attached to user successfully")
		attachPolicyToUser(accessKeyID, secretAccessKey, region, *result.Policy.Arn, userName)
		return
	}
	if err != nil {
		fmt.Println("Unable to attach role policy to user", err)
		return
	}
	attachPolicyToUser(accessKeyID, secretAccessKey, region, *result2.Policy.Arn, userName)
	fmt.Println("Policy attached to user successfully")
	success = true
	return
}
