package controllers

import (
	"context"

	egressv1 "github.com/monzo/egress-operator/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/api/networking/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// +kubebuilder:rbac:groups=networking,resources=networkpolicies,verbs=get;create;patch

func (r *ExternalServiceReconciler) reconcileNetworkPolicy(ctx context.Context, req ctrl.Request, es *egressv1.ExternalService) error {
	desired := networkPolicy(es)
	if err := ctrl.SetControllerReference(es, desired, r.Scheme); err != nil {
		return err
	}
	np := &v1.NetworkPolicy{}
	if err := r.Get(ctx, req.NamespacedName, np); err != nil {
		if apierrs.IsNotFound(err) {
			return r.Client.Create(ctx, desired)
		}
		return err
	}

	patched := np.DeepCopy()
	patched.Spec = desired.Spec

	return r.Client.Patch(ctx, patched, client.MergeFrom(np))
}

func networkPolicyPorts(es *egressv1.ExternalService) (ports []networkingv1.NetworkPolicyPort) {
	for _, port := range es.Spec.Ports {
		p := intstr.FromInt(int(port.Port))

		ports = append(ports, networkingv1.NetworkPolicyPort{
			Protocol: port.Protocol,
			Port:     &p,
		})
	}

	return
}

func networkPolicy(es *egressv1.ExternalService) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      es.Name,
			Namespace: namespace,
			Labels:    labels(es),
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: *metav1.SetAsLabelSelector(labelsToSelect(es)),
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"egress.monzo.com/allowed-gateway": es.Name,
								},
							},
							// Allow all namespaces
							NamespaceSelector: &metav1.LabelSelector{},
						},
					},
					Ports: networkPolicyPorts(es),
				},
			},
		},
	}
}
