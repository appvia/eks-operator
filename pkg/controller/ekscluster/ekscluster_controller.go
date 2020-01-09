package ekscluster

import (
	"context"
	"log"
	"time"

	awsv1alpha1 "github.com/appvia/eks-operator/pkg/apis/aws/v1alpha1"
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"
	"github.com/aws/aws-sdk-go/aws"
	eks "github.com/aws/aws-sdk-go/service/eks"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var logger = logf.Log.WithName("controller_ekscluster")

// Add creates a new EKSCluster Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileEKSCluster{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ekscluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource EKSCluster
	err = c.Watch(&source.Kind{Type: &awsv1alpha1.EKSCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner EKSCluster
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &awsv1alpha1.EKSCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileEKSCluster implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileEKSCluster{}

// ReconcileEKSCluster reconciles a EKSCluster object
type ReconcileEKSCluster struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileEKSCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := logger.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling EKSCluster")

	// Fetch the EKSCluster instance
	cluster := &awsv1alpha1.EKSCluster{}

	if err := r.client.Get(context.TODO(), request.NamespacedName, cluster); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	reqLogger.Info("Found AWSCluster CR")

	credentials := &awsv1alpha1.AWSCredential{}

	reference := types.NamespacedName{
		Namespace: cluster.Spec.Use.Namespace,
		Name:      cluster.Spec.Use.Name,
	}

	ctx := context.Background()

	err := r.client.Get(ctx, reference, credentials)

	if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Found AWSCredential CR")

	sesh, err := GetAWSSession(credentials, cluster.Spec.Region)

	svc, err := GetEKSService(sesh)

	reqLogger.Info("Checking cluster existence")
	exists, err := EKSClusterExists(svc, cluster.Spec.Name)

	if err != nil {
		return reconcile.Result{}, err
	}

	if exists {
		reqLogger.Info("Cluster exists:" + cluster.Spec.Name)
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Creating cluster:" + cluster.Spec.Name)

	// Cluster doesnt exist, create it
	_, err = CreateEKSCluster(svc, &eks.CreateClusterInput{
		Name:    aws.String(cluster.Spec.Name),
		RoleArn: aws.String(cluster.Spec.RoleArn),
		Version: aws.String(cluster.Spec.Version),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SecurityGroupIds: aws.StringSlice(cluster.Spec.SecurityGroupIds),
			SubnetIds:        aws.StringSlice(cluster.Spec.SubnetIds),
		},
	})

	if err != nil {
		return reconcile.Result{}, err
	}

	// Set status to pending
	cluster.Status.Status = core.PendingStatus

	if err := r.client.Status().Update(ctx, cluster); err != nil {
		logger.Error(err, "failed to update the resource status")
		return reconcile.Result{}, err
	}

	// Wait for it to become ACTIVE
	for {
		log.Println("Checking the status of cluster:", cluster.Spec.Name)

		status, err := GetEKSClusterStatus(svc, &eks.DescribeClusterInput{
			Name: aws.String(cluster.Spec.Name),
		})

		if err != nil {
			return reconcile.Result{}, err
		}

		if status == "ACTIVE" {
			log.Println("Cluster active:", cluster.Spec.Name)
			// Set status to success
			cluster.Status.Status = core.SuccessStatus

			if err := r.client.Status().Update(ctx, cluster); err != nil {
				logger.Error(err, "failed to update the resource status")
				return reconcile.Result{}, err
			}
			break
		}
		if status == "ERROR" {
			log.Println("Cluster has ERROR status:", cluster.Spec.Name)
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}
	return reconcile.Result{}, nil
}
