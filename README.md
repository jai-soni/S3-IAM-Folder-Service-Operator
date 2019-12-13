# AWS S3 Folder Service Operator


The AWS Folder Service Operator provides an extension to AWS S3 bucket service
that automates creation of self healing S3 folder for a given user.
A Chart release is described through a Kubernetes custom resource named team2-kubeop-helmchart. 


## AWS S3 Folder Service Operator features

* Creates IAM User
* Creates Bucket with the folder named after the user's username
* Create a bucket policy restricting access of the user to their own folder and binding user and bucker folder
* Secret of the IAM user is stored in user defined IAM secret name
* Self healing for deleted IAM secret (Event driven)
* Self AWS resources (Time driven)



## Expected Custome Resource (CR) Definition
```yaml
apiVersion: app.s3folder.com/v1alpha1
kind: FolderService
metadata:
  name: example-folderservice
  namespace: namespace
spec:
  userName: user
  userSecret:
    name: user-secret
  platformSecrets:
    aws:
      credentials:
        name: iam-secret
```



## Commands for setup
```bash
minikube start
minikube dashboard

kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml

kubectl apply -f secrets/mysecrets.yaml

kubectl apply -f deploy/crds/app.s3folder.com_folderservices_crd.yaml
kubectl apply -f deploy/crds/app.s3folder.com_v1alpha1_folderservice_cr.yaml

operator-sdk generate k8s && operator-sdk generate openapi

operator-sdk up local
```


## Developer information

| Name | NEU ID | Email Address | Github username |
| --- | --- | --- | --- |
| Jai Soni| 001822913|soni.j@husky.neu.edu | jai-soni |
| Riddhi Kakadiya| 001811354 | kamlesh.r@husky.neu.edu | riddhiKakadiya |
| Sreerag Mandakathil Sreenath| 001838559| mandakathil.s@husky.neu.edu| sreeragsreenath |
| Vivek Dalal| 001430934 | dalal.vi@husky.neu.edu | vivdalal |

Your feedback is always welcome!
