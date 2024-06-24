/*
Copyright 2024 Timoni Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1 "github.com/stefanprodan/timoni/api/v1alpha1"
	"github.com/stefanprodan/timoni/internal/controllers/api/v1alpha1"
	"github.com/stefanprodan/timoni/internal/reconciler"

	"github.com/stefanprodan/timoni/internal/engine"
	"github.com/stefanprodan/timoni/internal/engine/fetcher"
	"github.com/stefanprodan/timoni/internal/runtime"
)

type applyOpts struct {
	scheme             *kruntime.Scheme
	cacheDir           string
	overwriteOwnership bool
	force              bool
	wait               bool
}

func apply(ctx context.Context, instance *v1alpha1.ClusterProject, opts applyOpts, timeout time.Duration) error {
	kubeconfigArgs := genericclioptions.NewConfigFlags(false)

	log := log.FromContext(ctx)

	tmpDir, err := os.MkdirTemp("", apiv1.FieldManager)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	ctxPull, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	f, err := fetcher.New(ctxPull, fetcher.Options{
		Source:      instance.Spec.Source.Repository,
		Version:     instance.Spec.Source.Tag,
		Destination: tmpDir,
		CacheDir:    opts.cacheDir,
		// Creds:        applyArgs.creds.String(),
		// Insecure:     false,
		// DefaultLocal: true,
	})
	if err != nil {
		return err
	}
	mod, err := f.Fetch()
	if err != nil {
		return err
	}

	cuectx := cuecontext.New()
	builder := engine.NewModuleBuilder(
		cuectx,
		instance.Name,
		instance.Namespace,
		f.GetModuleRoot(),
		"main",
	)

	if err := builder.WriteSchemaFile(); err != nil {
		return err
	}

	mod.Name, err = builder.GetModuleName()
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("using module %s version %s", mod.Name, mod.Version))

	/*
		if len(applyArgs.valuesFiles) > 0 {
			valuesCue, err := convertToCue(cmd, applyArgs.valuesFiles)
			if err != nil {
				return err
			}
			err = builder.MergeValuesFile(valuesCue)
			if err != nil {
				return err
			}
		}
	*/

	kubeVersion, err := runtime.ServerVersion(kubeconfigArgs)
	if err != nil {
		return err
	}

	builder.SetVersionInfo(mod.Version, kubeVersion)

	buildResult, err := builder.Build()
	if err != nil {
		return describeErr(f.GetModuleRoot(), "build failed", err)
	}

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	bundleInstance := &engine.BundleInstance{
		Name:      instance.Name,
		Namespace: instance.Namespace,
		Module:    *mod,
		Bundle:    "",
	}

	r := reconciler.NewInteractiveReconciler(log,
		&reconciler.CommonOptions{
			Dir:                tmpDir,
			Wait:               opts.wait,
			Force:              opts.force,
			OverwriteOwnership: opts.overwriteOwnership,
		},
		&reconciler.InteractiveOptions{},
		timeout,
	)
	if err := r.Init(ctx, builder, buildResult, bundleInstance, kubeconfigArgs); err != nil {
		if errors.Is(err, &reconciler.InstanceOwnershipConflictErr{}) {
			return fmt.Errorf("%s %s", err, "apply with \"--overwrite-ownership\" to gain instance ownership.")
		}
		return err
	}
	return r.ApplyInstance(ctx, log,
		builder,
		buildResult,
	)
}

func describeErr(moduleRoot, description string, err error) error {
	return fmt.Errorf("%s:\n%s", description, errors.Details(err, &errors.Config{
		Cwd: moduleRoot,
	}))
}
