package api

//go:generate ${GOPATH}/bin/controller-gen object:headerFile="../../../api/license.go.txt" paths=./...
//go:generate ${GOPATH}/bin/controller-gen crd paths=./... output:dir=./crds
