package main

import "k8s/pkg/controller"

func main() {
	m := controller.NewManager()
	m.Start()
}
