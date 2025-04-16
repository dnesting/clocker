package clocker_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/dnesting/clocker"
)

func TestClocker_AddAndToString(t *testing.T) {
	c := clocker.Clocker{}
	ts := clocker.Timestamp{T: time.Now(), Annotation: []string{"test annotation"}}
	idx := c.Add(ts)

	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}

	output := c.ToString(idx)
	if output == "" {
		t.Errorf("expected non-empty output, got empty string")
	}
}

func TestClocker_Annotate(t *testing.T) {
	c := clocker.Clocker{}
	ts := clocker.Timestamp{T: time.Now()}
	c.Add(ts)

	idx, err := c.Annotate(0, "new annotation")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}

	if len(c.Timestamps[0].Annotation) != 1 || c.Timestamps[0].Annotation[0] != "new annotation" {
		t.Errorf("annotation not added correctly")
	}
}

func TestClocker_WriteTo(t *testing.T) {
	c := clocker.Clocker{}
	c.Add(clocker.Timestamp{T: time.Now(), Annotation: []string{"annotation 1"}})
	c.Add(clocker.Timestamp{T: time.Now(), Annotation: []string{"annotation 2"}})

	var buf bytes.Buffer
	written, err := c.WriteTo(&buf)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if written == 0 {
		t.Errorf("expected non-zero bytes written, got %d", written)
	}

	if buf.Len() == 0 {
		t.Errorf("expected non-empty buffer, got empty")
	}
}
