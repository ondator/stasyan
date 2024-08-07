package main

import (
	"context"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	namespace    string
	destroy_till time.Time
}

func main() {
	conf := readConf()
	destroy(conf)
}

func destroy(conf Config) {
	for checkTime(conf.destroy_till) {
		pods := getPods(conf.namespace)
		pod := chooseTarget(pods)
		kill(pod, conf.namespace)
		wait()
	}
}

func wait() {
	time.Sleep(5 * time.Minute)
}

func kill(pod string, namespace string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalln("can't obtain k8s conf:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can't create k8s api client: ", err)
	}
	err = clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod, metav1.DeleteOptions{})
	if err != nil {
		log.Println("[WARN] can't kill pod ", pod, err)
	}
}

func chooseTarget(pods []string) string {
	return pods[rand.IntN(len(pods)-1)]
}

func getPods(namespace string) []string {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalln("can't obtain k8s conf:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("can't create k8s api client: ", err)
	}
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []string{}
	}

	pods := Map(podList.Items, func(p v1.Pod) string {
		return p.Name
	})

	return pods
}

func checkTime(till time.Time) bool {
	return till.After(time.Now())
}

func readConf() Config {
	ns, defn := os.LookupEnv("STASYAN_NAMESPACE")
	if !defn {
		ns = "default"
	}

	t, deft := os.LookupEnv("STASYAN_LIFETIME")
	min, err := strconv.Atoi(t)
	if !deft || err != nil {
		min = 60
	}

	return Config{namespace: ns, destroy_till: time.Now().Add(time.Minute * time.Duration(min))}
}

func Map[T interface{}, F interface{}](t []T, f func(T) F) []F {
	res := make([]F, len(t))
	for i, v := range t {
		res[i] = f(v)
	}
	return res
}
