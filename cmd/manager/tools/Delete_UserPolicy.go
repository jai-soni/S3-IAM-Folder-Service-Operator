package main

import (
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/iam"
)

func Delete_UserPolicy() {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1")},
    )
    svc := iam.New(sess)

    foundPolicy := false
    policyName := "sreerag_policy"
	policyArn := "arn:aws:iam::aws:policy/sreerag_policy"
	username := "sreerag"

    _, err = svc.DetachUserPolicy(&iam.DetachUserPolicyInput{
        PolicyArn: &policyArn,
        UserName:  &username,
    })

    if err != nil {
        fmt.Println("Unable to detach policy to user")
        return
    }
    fmt.Println("policy detached successfully")
}

func Delete_User(){
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1")},
    )
    svc := iam.New(sess)

    _, err = svc.DeleteUser(&iam.DeleteUserInput{
        UserName: &username,
	})
    // If the user does not exist than we will log an error.	
	if awserr, ok := err.(awserr.Error); ok && awserr.Code() == iam.ErrCodeNoSuchEntityException {
        fmt.Printf("User %s does not exist\n", &username)
        return
    } else if err != nil {
        fmt.Println("Error", err)
        return
    }
    fmt.Printf("User %s has been deleted\n", &username)
}

func main()
{
	Delete_UserPolicy()
	Delete_User()
}