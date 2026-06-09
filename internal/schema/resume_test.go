package schema

import (
	"encoding/json"
	"testing"
)

func TestResumeUnmarshalProjectExperienceArrayFields(t *testing.T) {
	data := []byte(`{
		"basics": {},
		"project_experiences": [
			{
				"name": "AResumeI",
				"role": "Backend",
				"start_date": "2026-01",
				"end_date": "2026-02",
				"description": "Resume parser",
				"highlights": ["parsed resume JSON"],
				"technologies": ["Go"],
				"links": ["https://example.com"]
			}
		],
		"education": []
	}`)

	var resume Resume
	if err := json.Unmarshal(data, &resume); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	got := resume.ProjectExperiences[0]
	if got.Links[0] != "https://example.com" {
		t.Fatalf("Links[0] = %q", got.Links[0])
	}
}
