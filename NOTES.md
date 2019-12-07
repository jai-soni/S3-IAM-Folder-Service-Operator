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

