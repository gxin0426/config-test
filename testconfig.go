// package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"
// 	"k8s.io/client-go/util/homedir"
// )

// func main() {
// 	// kubeConfig, err := CreateKubeConfig()
// 	//当程序以pod方式运行时，就直接走这里的逻辑
// 	kubeConfig, err := rest.InClusterConfig()
// 	if err != nil {
// 		fmt.Println("config err~~ :", err)
// 	}

// 	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
// 	if err != nil {
// 		fmt.Println("crate clinet err : ", err)
// 	}
// 	//获取pod资源
// 	pods, err := kubeClient.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
// 	fmt.Println(pods)
// }

// func PathExists(path string) (bool, error) {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if os.IsNotExist(err) {
// 		return false, nil
// 	}
// 	return false, err
// }

// func CreateKubeConfig() (*rest.Config, error) {
// 	kubeConfigPath := ""
// 	if home := homedir.HomeDir(); home != "" {
// 		kubeConfigPath = filepath.Join(home, ".kube", "config")
// 	}
// 	fileExist, err := PathExists(kubeConfigPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("justify kubeConfigPath exist err,err:%v", err)
// 	}
// 	//.kube/config文件存在，就使用文件
// 	//这里主要是本地测试
// 	if fileExist {
// 		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return config, nil
// 	} else {
// 		//当程序以pod方式运行时，就直接走这里的逻辑
// 		config, err := rest.InClusterConfig()
// 		if err != nil {
// 			return nil, err
// 		}
// 		return config, nil
// 	}
// }

package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("kube-system").Get(context.TODO(), "coredns-6d8c4cb4d-2xltp", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod coredns-6d8c4cb4d-2xltp not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found coredns-6d8c4cb4d-2xltp pod in default namespace\n")
		}

		time.Sleep(10 * time.Second)
	}
}
