package main

import (
	"fmt"
	ofo "github.com/open-feature/open-feature-operator/apis/core/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"sigs.k8s.io/yaml"
)

func main() {

	nFlags := 5000

	sch := runtime.NewScheme()
	_ = scheme.AddToScheme(sch)

	ofo.AddToScheme(sch)

	flags := ofo.FeatureFlagSpec{
		Flags: make(map[string]ofo.FlagSpec, nFlags),
	}

	for i := 0; i < nFlags; i++ {
		flagName := fmt.Sprintf("color-%d", i)
		flag := ofo.FlagSpec{
			State:          "ENABLED",
			DefaultVariant: "blue",
			Variants: []byte(`{
				"blue":  "0d507d",
				"red":   "c05543",
				"green": "2f5230"
			}`),
			Targeting: []byte(`{
								"if": [
								  {
									"sem_ver": [{"var": "version"}, ">", "0.1.0"]
								  },
								  "red", "green"
								]
							  }`),
		}

		flags.Flags[flagName] = flag

	}

	fsc := &ofo.FeatureFlagConfiguration{
		TypeMeta: v1.TypeMeta{
			APIVersion: "core.openfeature.dev/v1alpha2",
			Kind:       "FeatureFlagConfiguration",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "benchmark-flag-source-config",
			Namespace: "flagd-performance-test",
		},
		Spec: ofo.FeatureFlagConfigurationSpec{
			ServiceProvider: &ofo.FeatureFlagServiceProvider{
				Name: "flagd",
			},
			FeatureFlagSpec: flags,
		},
	}

	flagSourceConfigBytes, err := yaml.Marshal(fsc)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	err = os.WriteFile("manifests/flag-source-config.yaml", flagSourceConfigBytes, 0644)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}
