//
// Copyright 2021-2022 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitops

import (
	"reflect"
	"testing"

	routev1 "github.com/openshift/api/route/v1"
	appstudiov1alpha1 "github.com/redhat-appstudio/application-service/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestGenerateDeployment(t *testing.T) {
	componentName := "test-component"
	namespace := "test-namespace"
	replicas := int32(1)
	otherReplicas := int32(3)
	labels := map[string]string{
		"component": componentName,
	}

	tests := []struct {
		name           string
		component      appstudiov1alpha1.Component
		wantDeployment appsv1.Deployment
	}{
		{
			name: "Simple component, no optional fields set",
			component: appstudiov1alpha1.Component{
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
			},
			wantDeployment: appsv1.Deployment{
				TypeMeta: v1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &v1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: labels,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            "container-image",
									ImagePullPolicy: corev1.PullAlways,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Component, optional fields set",
			component: appstudiov1alpha1.Component{
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "test-component",
					Application:   "test-application",
					Replicas:      3,
					TargetPort:    5000,
					Build: appstudiov1alpha1.Build{
						ContainerImage: "quay.io/test/test-image:latest",
					},
					Env: []corev1.EnvVar{
						{
							Name:  "test",
							Value: "value",
						},
					},
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("2M"),
							corev1.ResourceMemory: resource.MustParse("1Gi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("1M"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},
			wantDeployment: appsv1.Deployment{
				TypeMeta: v1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &otherReplicas,
					Selector: &v1.LabelSelector{
						MatchLabels: labels,
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: labels,
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            "container-image",
									Image:           "quay.io/test/test-image:latest",
									ImagePullPolicy: corev1.PullAlways,
									Env: []corev1.EnvVar{
										{
											Name:  "test",
											Value: "value",
										},
									},
									Ports: []corev1.ContainerPort{
										{
											ContainerPort: int32(5000),
										},
									},
									ReadinessProbe: &corev1.Probe{
										InitialDelaySeconds: 10,
										PeriodSeconds:       10,
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.FromInt(5000),
											},
										},
									},
									LivenessProbe: &corev1.Probe{
										InitialDelaySeconds: 10,
										PeriodSeconds:       10,
										Handler: corev1.Handler{
											HTTPGet: &corev1.HTTPGetAction{
												Port: intstr.FromInt(5000),
												Path: "/",
											},
										},
									},
									Resources: corev1.ResourceRequirements{
										Limits: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("2M"),
											corev1.ResourceMemory: resource.MustParse("1Gi"),
										},
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("1M"),
											corev1.ResourceMemory: resource.MustParse("256Mi"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generatedDeployment := generateDeployment(tt.component)

			if !reflect.DeepEqual(*generatedDeployment, tt.wantDeployment) {
				t.Errorf("TestGenerateDeployment() error: expected %v got %v", tt.wantDeployment, generatedDeployment)
			}
		})
	}
}

func TestGenerateService(t *testing.T) {
	componentName := "test-component"
	namespace := "test-namespace"
	labels := map[string]string{
		"component": componentName,
	}

	tests := []struct {
		name        string
		component   appstudiov1alpha1.Component
		wantService corev1.Service
	}{
		{
			name: "Simple component object",
			component: appstudiov1alpha1.Component{
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "test-component",
					Application:   "test-application",
					TargetPort:    5000,
				},
			},
			wantService: corev1.Service{
				TypeMeta: v1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Service",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Selector: labels,
					Ports: []corev1.ServicePort{
						{
							Port:       int32(5000),
							TargetPort: intstr.FromInt(5000),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generatedService := generateService(tt.component)

			if !reflect.DeepEqual(*generatedService, tt.wantService) {
				t.Errorf("TestGenerateService() error: expected %v got %v", tt.wantService, generatedService)
			}
		})
	}
}

func TestGenerateRoute(t *testing.T) {
	componentName := "test-component"
	namespace := "test-namespace"
	labels := map[string]string{
		"component": componentName,
	}
	weight := int32(100)

	tests := []struct {
		name      string
		component appstudiov1alpha1.Component
		wantRoute routev1.Route
	}{
		{
			name: "Simple component object",
			component: appstudiov1alpha1.Component{
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "test-component",
					Application:   "test-application",
					TargetPort:    5000,
				},
			},
			wantRoute: routev1.Route{
				TypeMeta: v1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Route",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
					Labels:    labels,
				},
				Spec: routev1.RouteSpec{
					Port: &routev1.RoutePort{
						TargetPort: intstr.FromInt(5000),
					},
					TLS: &routev1.TLSConfig{
						InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
						Termination:                   routev1.TLSTerminationEdge,
					},
					To: routev1.RouteTargetReference{
						Kind:   "Service",
						Name:   componentName,
						Weight: &weight,
					},
				},
			},
		},
		{
			name: "Component object with route/hostname set",
			component: appstudiov1alpha1.Component{
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "test-component",
					Application:   "test-application",
					TargetPort:    5000,
					Route:         "example.com",
				},
			},
			wantRoute: routev1.Route{
				TypeMeta: v1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Route",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      componentName,
					Namespace: namespace,
					Labels:    labels,
				},
				Spec: routev1.RouteSpec{
					Host: "example.com",
					Port: &routev1.RoutePort{
						TargetPort: intstr.FromInt(5000),
					},
					TLS: &routev1.TLSConfig{
						InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
						Termination:                   routev1.TLSTerminationEdge,
					},
					To: routev1.RouteTargetReference{
						Kind:   "Service",
						Name:   componentName,
						Weight: &weight,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generatedRoute := generateRoute(tt.component)

			if !reflect.DeepEqual(*generatedRoute, tt.wantRoute) {
				t.Errorf("TestGenerateRoute() error: expected %v got %v", tt.wantRoute, generatedRoute)
			}
		})
	}
}
