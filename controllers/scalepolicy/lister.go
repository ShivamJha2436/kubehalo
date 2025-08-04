package scalepolicy

import (
	scalev1 "github.com/ShivamJha2436/kubehalo/api/v1"
	informers "github.com/ShivamJha2436/kubehalo/generated/informers/externalversions/kubehalo/v1"
	listers "github.com/ShivamJha2436/kubehalo/generated/listers/kubehalo/v1"
)

type ScalePolicyLister struct {
	lister listers.ScalePolicyLister
}

func NewScalePolicyLister(informer informers.ScalePolicyInformer) *ScalePolicyLister {
	return &ScalePolicyLister{
		lister: informer.Lister(),
	}
}

func (s *ScalePolicyLister) ListAll() ([]*scalev1.ScalePolicy, error) {
	return s.lister.List(nil)
}

func (s *ScalePolicyLister) Get(namespace, name string) (*scalev1.ScalePolicy, error) {
	return s.lister.ScalePolicies(namespace).Get(name)
}
