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
	"github.com/fluxcd/pkg/ssa"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1 "github.com/stefanprodan/timoni/api/v1alpha1"
	"github.com/stefanprodan/timoni/internal/controllers/api/v1alpha1"

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

	ctxPull, cancel := context.WithTimeout(context.Background(), timeout)
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

	finalValues, err := builder.GetDefaultValues()
	if err != nil {
		return fmt.Errorf("failed to extract values: %w", err)
	}

	applySets, err := builder.GetApplySets(buildResult)
	if err != nil {
		return fmt.Errorf("failed to extract objects: %w", err)
	}

	var objects []*unstructured.Unstructured
	for _, set := range applySets {
		objects = append(objects, set.Objects...)
	}

	rm, err := runtime.NewResourceManager(kubeconfigArgs)
	if err != nil {
		return err
	}

	rm.SetOwnerLabels(objects, instance.Name, instance.Namespace)

	for i := range objects {
		if objects[i].GetNamespace() == "" && instance.GetNamespace() != "" {
			// avoid the case of "cluster-scoped resource must not have a namespace-scoped owner" errors
			// TODO: this means cluster-scoped resources will need custom deletion, perhaps it forbids
			// owner-based deletion mechanism; assume it migth suffice for the MVP
			continue
		}
		if err := controllerutil.SetControllerReference(instance, objects[i], opts.scheme); err != nil {
			return fmt.Errorf("cannot set controller ownership (object=%s): %w", objects[i], err)
		}
	}

	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	exists := false
	sm := runtime.NewStorageManager(rm)
	storedInstance, err := sm.Get(ctx, instance.Name, instance.Namespace)
	if err == nil {
		exists = true
	}

	nsExists, err := sm.NamespaceExists(ctx, instance.Namespace)
	if err != nil {
		return fmt.Errorf("instance init failed: %w", err)
	}

	if !opts.overwriteOwnership && exists {
		err = instanceOwnershipConflicts(*storedInstance)
		if err != nil {
			return err
		}
	}

	im := runtime.NewInstanceManager(instance.Name, instance.Namespace, finalValues, *mod)

	if err := im.AddObjects(objects); err != nil {
		return fmt.Errorf("adding objects to instance failed: %w", err)
	}

	staleObjects, err := sm.GetStaleObjects(ctx, &im.Instance)
	if err != nil {
		return fmt.Errorf("getting stale objects failed: %w", err)
	}

	/*
		if applyArgs.dryrun || applyArgs.diff {
			if !nsExists {
				log.Info(colorizeJoin(colorizeNamespaceFromArgs(), ssa.CreatedAction, dryRunServer))
			}
			return instanceDryRunDiff(logr.NewContext(ctx, log), rm, objects, staleObjects, nsExists, tmpDir, applyArgs.diff)
		}
	*/

	if !exists {
		log.Info(fmt.Sprintf("installing %s in namespace %s", instance.Name, instance.Namespace))

		if err := sm.Apply(ctx, &im.Instance, true); err != nil {
			return fmt.Errorf("instance init failed: %w", err)
		}

		if !nsExists {
			log.Info("Namespace/"+instance.Namespace, ssa.CreatedAction)
		}
	} else {
		log.Info(fmt.Sprintf("upgrading %s in namespace %s", instance.Name, instance.Namespace))
	}

	applyOpts := runtime.ApplyOptions(opts.force, timeout)
	applyOpts.WaitInterval = 5 * time.Second

	waitOptions := ssa.WaitOptions{
		Interval: applyOpts.WaitInterval,
		Timeout:  timeout,
		FailFast: true,
	}

	for _, set := range applySets {
		if len(applySets) > 1 {
			log.Info(fmt.Sprintf("applying %s", set.Name))
		}

		cs, err := rm.ApplyAllStaged(ctx, set.Objects, applyOpts)
		if err != nil {
			return err
		}
		for _, change := range cs.Entries {
			log.Info(change.String())
		}

		if opts.wait {
			// spin := StartSpinner(fmt.Sprintf("waiting for %v resource(s) to become ready...", len(set.Objects)))
			err = rm.Wait(set.Objects, waitOptions)
			// spin.Stop()
			if err != nil {
				return err
			}
			log.Info("resources are ready")
		}

	}

	if images, err := builder.GetContainerImages(buildResult); err == nil {
		im.Instance.Images = images
	}

	if err := sm.Apply(ctx, &im.Instance, true); err != nil {
		return fmt.Errorf("storing instance failed: %w", err)
	}

	var deletedObjects []*unstructured.Unstructured
	if len(staleObjects) > 0 {
		deleteOpts := runtime.DeleteOptions(instance.Name, instance.Namespace)
		changeSet, err := rm.DeleteAll(ctx, staleObjects, deleteOpts)
		if err != nil {
			return fmt.Errorf("pruning objects failed: %w", err)
		}
		deletedObjects = runtime.SelectObjectsFromSet(changeSet, ssa.DeletedAction)
		for _, change := range changeSet.Entries {
			log.Info(change.String())
		}
	}

	if opts.wait {
		if len(deletedObjects) > 0 {
			// spin := StartSpinner(fmt.Sprintf("waiting for %v resource(s) to be finalized...", len(deletedObjects)))
			err = rm.WaitForTermination(deletedObjects, waitOptions)
			// spin.Stop()
			if err != nil {
				return fmt.Errorf("waiting for termination failed: %w", err)
			}

			log.Info("all resources are ready")
		}
	}

	return nil
}

func instanceOwnershipConflicts(instance apiv1.Instance) error {
	if currentOwnerBundle := instance.Labels[apiv1.BundleNameLabelKey]; currentOwnerBundle != "" {
		return fmt.Errorf("instance ownership conflict encountered. Apply with \"--overwrite-ownership\" to gain instance ownership. Conflict: instance \"%s\" exists and is managed by bundle \"%s\"", instance.Name, currentOwnerBundle)
	}
	return nil
}

func describeErr(moduleRoot, description string, err error) error {
	return fmt.Errorf("%s:\n%s", description, errors.Details(err, &errors.Config{
		Cwd: moduleRoot,
	}))
}
