module github.com/iximiuz/conman

go 1.16

require (
	github.com/emicklei/go-restful v2.16.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2
	github.com/opencontainers/runtime-tools v0.9.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023
	google.golang.org/grpc v1.38.0
	k8s.io/api v0.22.2 // indirect
	k8s.io/client-go v0.22.2
	k8s.io/cri-api v0.22.2
	k8s.io/kubernetes v1.22.2
)

replace (
	k8s.io/api => k8s.io/kubernetes/staging/src/k8s.io/api v0.0.0-20210915213132-8b5a19147530
	k8s.io/apiextensions-apiserver => k8s.io/kubernetes/staging/src/k8s.io/apiextensions-apiserver v0.0.0-20210915213132-8b5a19147530
	k8s.io/apimachinery => k8s.io/kubernetes/staging/src/k8s.io/apimachinery v0.0.0-20210915213132-8b5a19147530
	k8s.io/apiserver => k8s.io/kubernetes/staging/src/k8s.io/apiserver v0.0.0-20210915213132-8b5a19147530
	k8s.io/cli-runtime => k8s.io/kubernetes/staging/src/k8s.io/cli-runtime v0.0.0-20210915213132-8b5a19147530
	k8s.io/client-go => k8s.io/kubernetes/staging/src/k8s.io/client-go v0.0.0-20210915213132-8b5a19147530
	k8s.io/cloud-provider => k8s.io/kubernetes/staging/src/k8s.io/cloud-provider v0.0.0-20210915213132-8b5a19147530
	k8s.io/cluster-bootstrap => k8s.io/kubernetes/staging/src/k8s.io/cluster-bootstrap v0.0.0-20210915213132-8b5a19147530
	k8s.io/code-generator => k8s.io/kubernetes/staging/src/k8s.io/code-generator v0.0.0-20210915213132-8b5a19147530
	k8s.io/component-base => k8s.io/kubernetes/staging/src/k8s.io/component-base v0.0.0-20210915213132-8b5a19147530
	k8s.io/component-helpers => k8s.io/kubernetes/staging/src/k8s.io/component-helpers v0.0.0-20210915213132-8b5a19147530
	k8s.io/controller-manager => k8s.io/kubernetes/staging/src/k8s.io/controller-manager v0.0.0-20210915213132-8b5a19147530
	k8s.io/cri-api => k8s.io/kubernetes/staging/src/k8s.io/cri-api v0.0.0-20210915213132-8b5a19147530
	k8s.io/csi-translation-lib => k8s.io/kubernetes/staging/src/k8s.io/csi-translation-lib v0.0.0-20210915213132-8b5a19147530
	k8s.io/kube-aggregator => k8s.io/kubernetes/staging/src/k8s.io/kube-aggregator v0.0.0-20210915213132-8b5a19147530
	k8s.io/kube-controller-manager => k8s.io/kubernetes/staging/src/k8s.io/kube-controller-manager v0.0.0-20210915213132-8b5a19147530
	k8s.io/kube-proxy => k8s.io/kubernetes/staging/src/k8s.io/kube-proxy v0.0.0-20210915213132-8b5a19147530
	k8s.io/kube-scheduler => k8s.io/kubernetes/staging/src/k8s.io/kube-scheduler v0.0.0-20210915213132-8b5a19147530
	k8s.io/kubectl => k8s.io/kubernetes/staging/src/k8s.io/kubectl v0.0.0-20210915213132-8b5a19147530
	k8s.io/kubelet => k8s.io/kubernetes/staging/src/k8s.io/kubelet v0.0.0-20210915213132-8b5a19147530
	k8s.io/kubernetes => k8s.io/kubernetes v1.22.2
	k8s.io/legacy-cloud-providers => k8s.io/kubernetes/staging/src/k8s.io/legacy-cloud-providers v0.0.0-20210915213132-8b5a19147530
	k8s.io/metrics => k8s.io/kubernetes/staging/src/k8s.io/metrics v0.0.0-20210915213132-8b5a19147530
	k8s.io/mount-utils => k8s.io/kubernetes/staging/src/k8s.io/mount-utils v0.0.0-20210915213132-8b5a19147530
	k8s.io/pod-security-admission => k8s.io/kubernetes/staging/src/k8s.io/pod-security-admission v0.0.0-20210915213132-8b5a19147530
	k8s.io/sample-apiserver => k8s.io/kubernetes/staging/src/k8s.io/sample-apiserver v0.0.0-20210915213132-8b5a19147530
)
