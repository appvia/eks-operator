package eksnodegroup

import (
	"context"
	"time"

	awsv1alpha1 "github.com/appvia/eks-operator/pkg/apis/aws/v1alpha1"
	core "github.com/appvia/hub-apis/pkg/apis/core/v1"
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

var logger = logf.Log.WithName("controller_eksnodegroup")

// Add creates a new EKSNodeGroup Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileEKSNodeGroup{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("eksnodegroup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource EKSNodeGroup
	err = c.Watch(&source.Kind{Type: &awsv1alpha1.EKSNodeGroup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner EKSNodeGroup
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &awsv1alpha1.EKSNodeGroup{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileEKSNodeGroup implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileEKSNodeGroup{}

// ReconcileEKSNodeGroup reconciles a EKSNodeGroup object
type ReconcileEKSNodeGroup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileEKSNodeGroup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := logger.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling EKSNodeGroup")

	// Fetch the EKSNodeGroup instance
	nodegroup := &awsv1alpha1.EKSNodeGroup{}

	if err := r.client.Get(context.TODO(), request.NamespacedName, nodegroup); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	reqLogger.Info("Found AWSNodeGroup CR")

	credentials := &awsv1alpha1.AWSCredential{}

	reference := types.NamespacedName{
		Namespace: nodegroup.Spec.Use.Namespace,
		Name:      nodegroup.Spec.Use.Name,
	}

	ctx := context.Background()

	err := r.client.Get(ctx, reference, credentials)

	if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Found AWSCredential CR")

	sesh, err := GetAWSSession(credentials, nodegroup.Spec.Region)

	svc, err := GetEKSService(sesh)

	existingNodeGroup, err := DescribeEKSNodeGroup(svc, &eks.DescribeNodegroupInput{
		ClusterName:   nodegroup.Spec.ClusterName,
		NodegroupName: nodegroup.Spec.NodegroupName,
	})

	if existingNodeGroup.Nodegroup {
		reqLogger.Info("Nodegroup exists")
		return reconcile.Result{}, nil
	}

	// Set status to pending
	nodegroup.Status.Status = core.PendingStatus

	if err := r.client.Status().Update(ctx, nodegroup); err != nil {
		logger.Error(err, "failed to update the resource status")
		return reconcile.Result{}, err
	}

	// Creat node group
	output, err := CreateEKSNodeGroup(svc, &eks.CreateNodegroupInput{
		ClusterName:   nodegroup.Spec.ClusterName,
		NodeRole:      nodegroup.Spec.NodeRole,
		NodegroupName: nodegroup.Spec.NodegroupName,
		Subnets:       nodegroup.Spec.Subnets,
	})

	if err != nil {
		return reconcile.Result{}, err
	}

	// Wait for node group to become ACTIVE
	for {
		reqLogger.Info("Checking the status of the node group:" + nodegroup.Spec.NodegroupName)

		nodestatus, err := GetEKSNodeGroupStatus(svc, &eks.DescribeNodegroupInput{
			ClusterName:   aws.String(nodegroup.Spec.ClusterName),
			NodegroupName: aws.String(nodegroup.Spec.NodegroupName),
		})

		if err != nil {
			return reconcile.Result{}, err
		}

		if nodestatus == "ACTIVE" {
			reqLogger.Info("Nodegroup active:" + nodegroup.Spec.NodegroupName)
			// Set status to success
			nodegroup.Status.Status = core.SuccessStatus

			if err := r.client.Status().Update(ctx, nodegroup); err != nil {
				logger.Error(err, "failed to update the resource status")
				return reconcile.Result{}, err
			}
			break
		}
		if nodestatus == "ERROR" {
			reqLogger.Info("Node group has ERROR status:" + nodegroup.Spec.NodegroupName)
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}

	return reconcile.Result{}, nil
}
