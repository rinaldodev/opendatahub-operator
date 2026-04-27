package precondition

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
)

func MonitorCRD(gvk schema.GroupVersionKind, opts ...Option) PreCondition {
	return MonitorCRDs([]schema.GroupVersionKind{gvk}, opts...)
}

// All CRDs are checked and absent ones are reported together in a single failure message.
// The first API error encountered is returned immediately.
func MonitorCRDs(gvks []schema.GroupVersionKind, opts ...Option) PreCondition {
	return newPreCondition(func(ctx context.Context, rr *types.ReconciliationRequest) (checkResult, error) {
		if len(gvks) == 0 {
			return checkResult{}, errors.New("MonitorCRDs called with empty GVK list")
		}

		var missing []string

		for _, gvk := range gvks {
			has, err := cluster.HasCRD(ctx, rr.Client, gvk)
			if err != nil {
				return checkResult{}, fmt.Errorf("%s: failed to check CRD presence: %w", gvk.Kind, err)
			}

			if !has {
				missing = append(missing, gvk.Kind+": CRD not found")
			}
		}

		if len(missing) > 0 {
			return checkResult{pass: false, message: strings.Join(missing, "; ")}, nil
		}

		return checkResult{pass: true}, nil
	}, opts...)
}
