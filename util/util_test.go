package util

import (
	"testing"
)

func TestSimpleMerge(t *testing.T) {
	a := make(map[string]interface{})
	b := make(map[string]interface{})
	a["foo"] = "bar"
	b["baz"] = "glub"
	err := MergeMap(a, b)
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, ok := a["baz"]; !ok {
		t.Errorf("key 'baz' was not in destination map, but it should have been.")
	}
}

func TestMoreComplexMerge(t *testing.T) {
	a := make(map[string]interface{})
	b := make(map[string]interface{})
	a["foo"] = "bar"
	acpu := make(map[string]interface{})
	acpu["model"] = "Ж86"
	a["cpujunk"] = acpu
	b["baz"] = 4532
	bcpu := make(map[string]interface{})
	bcpu["flags"] = []string{"fpu", "vme", "de", "pse", "tsc", "msr", "pae", "mce", "cx8", "apic", "sep", "mtrr", "pge", "mca", "cmov", "pat", "pse36", "clfsh", "ds", "acpi", "mmx"}
	b["foo"] = "meek"
	bcpu["model"] = "x86-64"
	b["cpujunk"] = bcpu
	err := MergeMap(a, b)
	if err != nil {
		t.Errorf(err.Error())
	}
	if a["foo"] != "bar" {
		t.Errorf("a[\"foo\"] should have been 'bar', but instead it was '%s'", a["foo"])
	}
	if aflags, ok := a["cpujunk"].(map[string]interface{})["flags"].([]string); !ok {
		t.Errorf("a[\"cpujunk\"][\"flags\"] should have been a slice of strings, but it was '%v' :: %T", a["cpujunk"].(map[string]interface{})["flags"], a["cpujunk"].(map[string]interface{})["flags"])
	} else {
		if len(aflags) != 21 {
			t.Errorf("somehow the length of the flags slice in a wasy wrong, expected 21 got %d", len(aflags))
		}
	}
	if a["cpujunk"].(map[string]interface{})["model"] != "Ж86" {
		t.Errorf("Looks like a[\"model\"] got overwritten with '%s'", a["cpujunk"].(map[string]interface{})["model"])
	}
}
