package folderservice

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sreeragsreenath/team2-kubeop/cmd/manager/tools/aws_s3_custom"
	appv1alpha1 "github.com/sreeragsreenath/team2-kubeop/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	bucketName              = "BUCKET_NAME"
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
	// Watch for changes to secondary resource Pods and requeue the owner FolderService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
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

	// Define a new Pod object
	pod := newPodForCR(instance)
	var secretName = "iam-secret"
	var namespace = "default"
	var region = "us-east-1"
	var userName = "jai1"
	var accessKeyID, secretAccessKey, bucketName = getCredentialsAndBucketDetails(secretName, namespace, region)
	aws_s3_custom.CreateUserIfNotExist(accessKeyID, secretAccessKey, userName, region)
	aws_s3_custom.CreateFolderIfNotExist(accessKeyID, secretAccessKey, userName+"/", bucketName, region)
	aws_s3_custom.CreatePolicyIfNotExist(accessKeyID, secretAccessKey, userName+"/", bucketName, region, userName)

	// Set FolderService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *appv1alpha1.FolderService) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func getCredentialsAndBucketDetails(secretName, namespace, region string) (string, string, string) {
	//awsConfig := &aws.Config{Region: aws.String(region)}
	s := scheme.Scheme
	config, err := getClientConfig()
	if err != nil {
		fmt.Errorf("Unable to get client configuration", err)
	}

	var kubeClient = getKubeClientOrDie(config, s)
	if secretName != "" {
		secret := &corev1.Secret{}
		err := kubeClient.Get(context.TODO(),
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
		accessKeyID, ok := secret.Data[awsCredsSecretIDKey]
		if !ok {
			fmt.Errorf("AWS credentials secret %v did not contain key %v",
				secretName, awsCredsSecretIDKey)
		}

		secretAccessKey, ok := secret.Data[awsCredsSecretAccessKey]
		if !ok {
			fmt.Errorf("AWS credentials secret %v did not contain key %v",
				secretName, awsCredsSecretAccessKey)
		}

		bucketNameForFolder, ok := secret.Data[bucketName]
		if !ok {
			fmt.Errorf("Bucket Name %v",
				secretName, bucketName)
		}

		return strings.Trim(string(accessKeyID), "\n"),
			strings.Trim(string(secretAccessKey), "\n"),
			strings.Trim(string(bucketNameForFolder), "\n")

	}
	return "", "", ""
}

func getKubeClientOrDie(config *rest.Config, s *runtime.Scheme) client.Client {
	c, err := client.New(config, client.Options{Scheme: s})
	if err != nil {
		panic(err)
	}
	return c
}

func getClientConfig() (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", path.Join(os.Getenv("HOME"), ".kube/config"))
}
