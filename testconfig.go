package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// kubeConfig, err := CreateKubeConfig()
	//当程序以pod方式运行时，就直接走这里的逻辑
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("config err~~ :", err)
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Println("crate clinet err : ", err)
	}
	//获取pod资源
	pods, err := kubeClient.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	fmt.Println(pods)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateKubeConfig() (*rest.Config, error) {
	kubeConfigPath := ""
	if home := homedir.HomeDir(); home != "" {
		kubeConfigPath = filepath.Join(home, ".kube", "config")
	}
	fileExist, err := PathExists(kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("justify kubeConfigPath exist err,err:%v", err)
	}
	//.kube/config文件存在，就使用文件
	//这里主要是本地测试
	if fileExist {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
		return config, nil
	} else {
		//当程序以pod方式运行时，就直接走这里的逻辑
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}
}
