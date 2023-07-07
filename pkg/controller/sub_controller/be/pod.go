package be

import (
	"context"
	v1 "github.com/selectdb/doris-operator/api/v1"
	"github.com/selectdb/doris-operator/pkg/common/utils/resource"
	corev1 "k8s.io/api/core/v1"
	"strconv"
)

func (be *Controller) buildBEPodTemplateSpec(dcr *v1.DorisCluster) corev1.PodTemplateSpec {
	podTemplateSpec := resource.NewPodTemplateSpc(dcr, v1.Component_BE)
	var containers []corev1.Container
	containers = append(containers, podTemplateSpec.Spec.Containers...)
	beContainer := be.beContainer(dcr)
	containers = append(containers, beContainer)
	podTemplateSpec.Spec.Containers = containers
	return podTemplateSpec
}

func (be *Controller) beContainer(dcr *v1.DorisCluster) corev1.Container {
	c := resource.NewBaseMainContainer(dcr.Spec.BeSpec.BaseSpec, v1.Component_BE)
	config, _ := be.GetConfig(context.Background(), &dcr.Spec.BeSpec.ConfigMapInfo, dcr.Namespace)
	addr := v1.GetConfigFEAddrForAccess(dcr, v1.Component_BE)
	var feconfig map[string]interface{}
	if addr == "" {
		if dcr.Spec.BeSpec.ConfigMapInfo.ConfigMapName != "" && dcr.Spec.BeSpec.ConfigMapInfo.ResolveKey != "" {
			feconfig, _ = be.getFeConfig(context.Background(), &dcr.Spec.BeSpec.ConfigMapInfo, dcr.Namespace)
		}
		config[resource.QUERY_PORT] = strconv.FormatInt(int64(resource.GetPort(feconfig, resource.QUERY_PORT)), 10)
	}

	ports := resource.GetContainerPorts(config, v1.Component_BE)
	c.Name = "be"
	c.Ports = append(c.Ports, ports...)
	c.Env = append(c.Env, corev1.EnvVar{
		Name:  resource.ENV_FE_ADDR,
		Value: addr,
	})

	return c
}
