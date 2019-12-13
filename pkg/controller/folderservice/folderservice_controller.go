package folderservice

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/sreeragsreenath/team2-kubeop/cmd/manager/tools/aws_s3_custom"
	appv1alpha1 "github.com/sreeragsreenath/team2-kubeop/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	awsCredsSecretIDKey     = "AWS_ACCESS_KEY_ID"
	awsCredsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	bucketNameFromSecret    = "BUCKET_NAME"
)

var log = logf.Log.WithName("controller_folderservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new FolderService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFolderService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("folderservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource FolderService
	err = c.Watch(&source.Kind{Type: &appv1alpha1.FolderService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Secret and requeue the owner FolderService
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.FolderService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFolderService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFolderService{}

// ReconcileFolderService reconciles a FolderService object
type ReconcileFolderService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a FolderService object and makes changes based on the state read
// and what is in the FolderService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileFolderService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling FolderService")

	// Fetch the FolderService instance
	instance := &appv1alpha1.FolderService{}
	fmt.Println("Before> " + strconv.FormatBool(instance.Status.SetupComplete))

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue

			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var secretName = "iam-secret"
	var namespace = "default"
	var region = "us-east-1"
	var userName = instance.Spec.UserName
	var userSecret = instance.Spec.UserSecret.Name
	var accessKeyID = ""
	var secretAccessKey = ""
	var bucketName = ""

	fmt.Println(userName)

	if secretName != "" {
		secret := &corev1.Secret{}
		err := r.client.Get(context.TODO(),
			types.NamespacedName{
				Name:      secretName,
				Namespace: namespace,
			},
			secret)

		fmt.Print(err)
		// if err != nil {
		// 	//return nil, err
		// 	return nil, nil
		// }
		accessKeyID1, ok := secret.Data[awsCredsSecretIDKey]
		if !ok {
			fmt.Errorf("AWS credentials secret %v did not contain key %v",
				secretName, awsCredsSecretIDKey)
		}

		secretAccessKey1, ok := secret.Data[awsCredsSecretAccessKey]
		if !ok {
			fmt.Errorf("AWS credentials secret %v did not contain key %v",
				secretName, awsCredsSecretAccessKey)
		}

		bucketNameForFolder1, ok := secret.Data[bucketNameFromSecret]
		if !ok {
			fmt.Errorf("Bucket Name error %v",
				secretName, bucketNameFromSecret)
		}

		fmt.Println("Bucket Name : %v", bucketNameForFolder1)

		accessKeyID = strings.Trim(string(accessKeyID1), "\n")
		secretAccessKey = strings.Trim(string(secretAccessKey1), "\n")
		bucketName = strings.Trim(string(bucketNameForFolder1), "\n")
	}

	//var accessKeyID, secretAccessKey, bucketName = getCredentialsAndBucketDetails(secretName, namespace, region)
	var resultAwsAccessKey, resultAwsSecretAccessKey, _ = aws_s3_custom.CreateUserIfNotExist(accessKeyID, secretAccessKey, userName, region)
	aws_s3_custom.CreateFolderIfNotExist(accessKeyID, secretAccessKey, userName+"/", bucketName, region)
	aws_s3_custom.CreatePolicyIfNotExist(accessKeyID, secretAccessKey, userName+"/", bucketName, region, userName)

	status := appv1alpha1.FolderServiceStatus{
		SetupComplete: true,
	}
	iamSecret := newIAMSecretCR(instance, userSecret, resultAwsAccessKey, resultAwsSecretAccessKey)
	// setOwnerRef(desired)

	if err := controllerutil.SetControllerReference(instance, iamSecret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Check if this Pod already exists
	found := &corev1.Secret{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: iamSecret.Name, Namespace: iamSecret.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", iamSecret.Namespace, "Pod.Name", iamSecret.Name)
		err = r.client.Create(context.TODO(), iamSecret)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(instance.Status, status) {
		instance.Status = status
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "failed to update the podSet")
			return reconcile.Result{}, err
		}
	}
	//instance.Status.SetupComplete = true
	fmt.Println("After> " + strconv.FormatBool(instance.Status.SetupComplete))

	return reconcile.Result{}, nil
}

func setOwnerRef(secret *corev1.Secret) error {
	ownerRef := generateOwnerRef(secret)
	secret.SetOwnerReferences(ownerRef)
	return nil
}

func generateOwnerRef(secret *corev1.Secret) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(secret, schema.GroupVersionKind{
			Group:   appv1alpha1.SchemeGroupVersion.Group,
			Version: appv1alpha1.SchemeGroupVersion.Version,
		}),
	}
}

// newIAMSecretCR returns a busybox pod with the same name/namespace as the cr
func newIAMSecretCR(cr *appv1alpha1.FolderService, userSecret, resultAwsAccessKey, resultAwsSecretAccessKey string) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      userSecret,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			"username": resultAwsAccessKey,
			"password": resultAwsSecretAccessKey,
		},
	}
}
