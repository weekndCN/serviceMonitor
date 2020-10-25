module github.com/serviceMonitor/m

go 1.14

replace (
	github.com/serviceMonitor/m/alert => /serviceMonitor/alert
	github.com/serviceMonitor/m/dingtalk => /serviceMonitor/dingtalk
	github.com/serviceMonitor/m/k8s => /serviceMonitor/k8s
	k8s.io/api => k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5-beta.1
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/kubectl => k8s.io/kubectl v0.16.5
)

require (
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-00010101000000-000000000000
)
