# Commands to create api and controller (Do not run this)
```bash
operator-sdk new team2-kubeop --repo github.com/sreeragsreenath/team2-kubeop

operator-sdk add api --api-version=app.s3folder.com/v1alpha1 --kind=FolderService

operator-sdk add controller --api-version=app.s3folder.com/v1alpha1 --kind=FolderService
```

# Commands for setup
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


# Commands to create secrets
- Create a folder with name "secrets" at the root level of project

- Add the mysecrets.yaml file inside the folder

- Update the following values with their base64  representations

```bash
data:
    AWS_ACCESS_KEY_ID: ""
    AWS_SECRET_ACCESS_KEY: ""
    BUCKET_NAME: ""
```

- Run the following command to create the secret on the k8s cluster
```bash
kubectl apply -f secrets/mysecrets.yaml
```

- In order to delete the secret, run:
```bash
kubectl delete secret iam-secret
```


