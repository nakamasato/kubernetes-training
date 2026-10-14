package main

import (
	"context"
	"encoding/json"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestMutatingPatch(t *testing.T) {
	for _, annotations := range []string{"", `,"annotations":{"keep":"value","access":"old"}`} {
		raw := []byte(`{"apiVersion":"v1","kind":"Pod","metadata":{"name":"demo"` + annotations + `},"spec":{"customField":"preserved"}}`)
		response := mutate(context.Background(), admission.Request{AdmissionRequest: admissionv1.AdmissionRequest{Object: runtime.RawExtension{Raw: raw}}})
		if !response.Allowed {
			t.Fatalf("mutation rejected: %v", response.Result)
		}
		patchBytes, err := json.Marshal(response.Patches)
		if err != nil {
			t.Fatal(err)
		}
		patch, err := jsonpatch.DecodePatch(patchBytes)
		if err != nil {
			t.Fatal(err)
		}
		updated, err := patch.Apply(raw)
		if err != nil {
			t.Fatalf("invalid patch: %v", err)
		}
		obj := &unstructured.Unstructured{}
		if err := json.Unmarshal(updated, obj); err != nil {
			t.Fatal(err)
		}
		got := obj.GetAnnotations()
		if got["access"] != "granted" || got["reason"] != "not so secret" {
			t.Fatalf("annotations: %v", got)
		}
		if annotations != "" && got["keep"] != "value" {
			t.Fatal("existing annotation was lost")
		}
		if value, _, _ := unstructured.NestedString(obj.Object, "spec", "customField"); value != "preserved" {
			t.Fatal("unknown field was lost")
		}
	}
}

func TestMutatingMalformedObject(t *testing.T) {
	response := mutate(context.Background(), admission.Request{AdmissionRequest: admissionv1.AdmissionRequest{Object: runtime.RawExtension{Raw: []byte(`{`)}}})
	if response.Allowed || response.Result.Code != 400 {
		t.Fatalf("expected bad request, got %+v", response)
	}
}
