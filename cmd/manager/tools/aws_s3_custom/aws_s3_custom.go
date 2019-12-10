package aws_s3_custom

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// Create user Bucket
func Create(accessKeyID, secretAccessKey, region, bucketName string) {
	fmt.Println(">>>>>>>>>>>>>>>>")
	fmt.Println("Key %v", accessKeyID)
	fmt.Println("Secret %v", secretAccessKey)
	fmt.Println("BucketName %v", bucketName)
	fmt.Println("Region %v", region)
	fmt.Println(">>>>>>>>>>>>>>>>")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
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

	// var bucket = "test-s3-folder-upload"
	var filename = "riddhi/"
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	fmt.Printf("Successfully uploaded %q to %q\n", filename, bucketName)

}
