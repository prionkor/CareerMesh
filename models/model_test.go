package models

import (
	"encoding/json"
	"math"
	"testing"
)

func TestJobDescriptionAnalysisResultJSONSchema(t *testing.T) {
	result := JobDescriptionAnalysisResult{
		SchemaVersion: JobDescriptionAnalysisSchemaVersion,
		Requirements: []JobDescriptionAnalysisRequirement{{
			Name: "Kubernetes", Description: "Experience operating Kubernetes in production",
			Category:  JobDescriptionAnalysisRequirementCategorySkill,
			Priority:  JobDescriptionAnalysisRequirementPriorityRequired,
			Assertion: JobDescriptionAnalysisAssertionExplicit, Weight: 0.95,
			Evidence: []JobDescriptionAnalysisEvidence{{Text: "Experience operating Kubernetes in production"}},
		}},
		Responsibilities: []JobDescriptionAnalysisResponsibility{},
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	want := `{"schema_version":"1","requirements":[{"name":"Kubernetes","description":"Experience operating Kubernetes in production","category":"skill","priority":"required","assertion":"explicit","weight":0.95,"evidence":[{"text":"Experience operating Kubernetes in production"}]}],"responsibilities":[]}`
	if string(encoded) != want {
		t.Fatalf("JSON = %s, want %s", encoded, want)
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("validate result: %v", err)
	}
}

func TestJobDescriptionAnalysisResultRequiresArrays(t *testing.T) {
	result := JobDescriptionAnalysisResult{SchemaVersion: JobDescriptionAnalysisSchemaVersion}
	if err := result.Validate(); err == nil {
		t.Fatal("expected missing requirements array to fail validation")
	}

	result.Requirements = []JobDescriptionAnalysisRequirement{}
	if err := result.Validate(); err == nil {
		t.Fatal("expected missing responsibilities array to fail validation")
	}

	result.Responsibilities = []JobDescriptionAnalysisResponsibility{}
	if err := result.Validate(); err != nil {
		t.Fatalf("empty arrays should be valid: %v", err)
	}
}

func TestJobDescriptionAnalysisResultValidatesEnumsAndWeights(t *testing.T) {
	result := JobDescriptionAnalysisResult{
		SchemaVersion: JobDescriptionAnalysisSchemaVersion,
		Requirements: []JobDescriptionAnalysisRequirement{{
			Category:  JobDescriptionAnalysisRequirementCategorySkill,
			Priority:  JobDescriptionAnalysisRequirementPriorityRequired,
			Assertion: JobDescriptionAnalysisAssertionExplicit,
			Weight:    0, Evidence: []JobDescriptionAnalysisEvidence{},
		}},
		Responsibilities: []JobDescriptionAnalysisResponsibility{},
	}
	if err := result.Validate(); err != nil {
		t.Fatalf("zero weight should be valid: %v", err)
	}

	result.Requirements[0].Weight = 1
	if err := result.Validate(); err != nil {
		t.Fatalf("unit weight should be valid: %v", err)
	}

	result.Requirements[0].Weight = 1.01
	if err := result.Validate(); err == nil {
		t.Fatal("expected out-of-range weight to fail validation")
	}

	for _, weight := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		result.Requirements[0].Weight = weight
		if err := result.Validate(); err == nil {
			t.Errorf("expected non-finite weight %v to fail validation", weight)
		}
	}
}
