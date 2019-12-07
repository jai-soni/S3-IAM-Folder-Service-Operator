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

kubectl apply -f deploy/crds/app.s3folder.com_folderservices_crd.yaml
kubectl apply -f deploy/crds/app.s3folder.com_v1alpha1_folderservice_cr.yaml

operator-sdk up local

operator-sdk generate k8s && operator-sdk generate openapi
```


