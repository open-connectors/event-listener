package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Generate the Random 256 byte token
func generateToken() string {
	token := make([]byte, 256)
	_, err := rand.Read(token)
	if err != nil {
		fmt.Errorf("Found error while generating the secret token - %v", err)
		panic(err)
	}

	base64EncodedToken := base64.StdEncoding.EncodeToString(token)
	return base64EncodedToken
}

// get Role Object
func getRoleObject() *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: v1.ObjectMeta{
			Name:      "tekton-role", // Name of the Role resource
			Namespace: "default",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"tekton.dev/v1beta1"}, // add required API groups
				Resources: []string{"taskruns"},           // add required resources
				Verbs:     []string{"get", "list"},        // Add all the operations
			},
		},
	}
}

// get RoleBinding Object
func getRoleBindingObject() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      "tekton-role-binding", // name of the role binding Resource
			Namespace: "default",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     "tekton-role", // This should be the name of created role
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "tekton-sa", // Name of the service account
				Namespace: "default",   // Name of the Namespace
			},
		},
	}
}

func GetSecureClientSet() (*dynamic.DynamicClient, error) {
	// ClientSet from Inside
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Errorf("Fail to build the k8s config. Error - %s", err)
		return nil, err
	}
	// inorder to create the dynamic Client set
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Fail to create the k8s client set. Errorf - %s", err)
		return nil, err
	}
	// corecclientSet, err := corev1client.NewForConfig(config)
	// if err != nil {
	// 	return nil, err
	// }

	// Metadata for creating the ServiceAccount
	saRequest := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name: "tekton-sa",
		},
	}

	// Metadata for creating secret
	secretRequest := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name: "tekton-secret",
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": "tekton-sa",
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
		Data: map[string][]byte{
			"token": []byte(generateToken()),
		},
	}

	_, err = clientSet.CoreV1().ServiceAccounts("default").Create(context.TODO(), saRequest, v1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	_, err = clientSet.CoreV1().Secrets("default").Create(context.TODO(), secretRequest, v1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	// We need to define the rule for the ClientSet associated Service account
	// get the Role
	role := getRoleObject()

	// get the Role binding object
	roleBinding := getRoleBindingObject()

	// Context to consider
	ctx := context.Background()

	// Create the Role
	_, err = clientSet.RbacV1().Roles("default").Create(ctx, role, v1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("Found error, while creating the Role, err - %s", err)
	}

	// Create Role binding
	_, err = clientSet.RbacV1().RoleBindings("default").Create(ctx, roleBinding, v1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("Found error while creating the rolebinding. Error - %s", err)
	}

	// Now create the new client set with Service account information
	// add the service account info to the config
	config.Impersonate = rest.ImpersonationConfig{
		UserName: fmt.Sprintf("system:serviceaccount:%s:%s", "default", "tekton-sa"),
	}

	// return client set with respect to new config
	return dynamic.NewForConfig(config)
}
