// Test adapter only: constructs fixtures and calls the learner's parser.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	vpn "localvpn"
)

type fixture struct {
	ID        string `json:"id"`
	Milestone string `json:"milestoneId"`
	Purpose   string `json:"purpose"`
	Input     struct {
		Nil         bool         `json:"nil"`
		Length      int          `json:"length"`
		Capacity    int          `json:"capacity"`
		FirstByte   byte         `json:"firstByte"`
		TotalLength uint16       `json:"totalLength"`
		Source      [4]byte      `json:"source"`
		Destination [4]byte      `json:"destination"`
		Patches     map[int]byte `json:"patches"`
	} `json:"input"`
	Expected json.RawMessage `json:"expected"`
}

type result struct {
	ID       string          `json:"id"`
	Passed   bool            `json:"passed"`
	Expected json.RawMessage `json:"expected"`
	Actual   any             `json:"actual"`
	Error    string          `json:"error,omitempty"`
}

func run(f fixture) (r result) {
	r.ID, r.Expected = f.ID, f.Expected
	defer func() {
		if p := recover(); p != nil {
			r.Passed = false
			r.Actual = nil
			r.Error = fmt.Sprintf("panic: %v", p)
		}
	}()
	var packet []byte
	if !f.Input.Nil {
		capacity := f.Input.Capacity
		if capacity < f.Input.Length {
			capacity = f.Input.Length
		}
		packet = make([]byte, f.Input.Length, capacity)
		// Write only the fixture bytes that exist, including for short inputs.
		put := func(i int, b byte) {
			if i < len(packet) {
				packet[i] = b
			}
		}
		put(0, f.Input.FirstByte)
		put(2, byte(f.Input.TotalLength>>8))
		put(3, byte(f.Input.TotalLength))
		for i := 0; i < 4; i++ {
			put(12+i, f.Input.Source[i])
			put(16+i, f.Input.Destination[i])
		}
		for i, b := range f.Input.Patches {
			if i < 0 || i >= len(packet) {
				panic("fixture patch outside input")
			}
			packet[i] = b
		}
	}
	info, err := vpn.ParseIPv4(packet)
	if err != nil {
		r.Actual = map[string]any{"accepted": false}
	} else {
		r.Actual = map[string]any{
			"accepted": true, "headerLength": info.HeaderLength,
			"totalLength": info.TotalLength, "source": info.Source,
			"destination": info.Destination,
		}
	}
	actualJSON, err := json.Marshal(r.Actual)
	if err != nil {
		panic(err)
	}
	var actual, expected any
	if err := json.Unmarshal(actualJSON, &actual); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(f.Expected, &expected); err != nil {
		panic(err)
	}
	r.Passed = reflect.DeepEqual(actual, expected)
	return r
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	fmt.Println(`{"cases":[]}`)
	os.Exit(1)
}

func main() {
	path := os.Getenv("PI_BUILD_CASES")
	review := path != ""
	milestone := os.Getenv("PI_BUILD_MILESTONE")
	if !review {
		path = "cases.json"
		if len(os.Args) > 1 {
			milestone = os.Args[1]
		}
	}
	projectBytes, err := os.ReadFile("project.json")
	if err != nil {
		fail(err)
	}
	var project struct {
		Milestones []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"milestones"`
	}
	if err := json.Unmarshal(projectBytes, &project); err != nil {
		fail(err)
	}
	if milestone == "" {
		for _, m := range project.Milestones {
			if m.Status != "complete" && m.Status != "completed" {
				milestone = m.ID
				break
			}
		}
		if milestone == "" && len(project.Milestones) > 0 {
			milestone = project.Milestones[len(project.Milestones)-1].ID
		}
	}
	allowed := map[string]bool{}
	found := false
	for _, m := range project.Milestones {
		allowed[m.ID] = true
		if m.ID == milestone {
			found = true
			break
		}
	}
	if !found {
		fail(fmt.Errorf("unknown milestone %q", milestone))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fail(err)
	}
	var fixtures []fixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		fail(err)
	}
	report := struct {
		Cases []result `json:"cases"`
	}{Cases: []result{}}
	seen := map[string]bool{}
	allPassed := true
	for _, f := range fixtures {
		// A teacher snapshot is authoritative: never filter its cases.
		if !review && (!allowed[f.Milestone] || (f.Purpose != "development" && f.Purpose != "regression")) {
			continue
		}
		if seen[f.ID] {
			fail(fmt.Errorf("duplicate case ID %q", f.ID))
		}
		seen[f.ID] = true
		r := run(f)
		report.Cases = append(report.Cases, r)
		if !r.Passed {
			allPassed = false
		}
		fmt.Fprintf(os.Stderr, "%s: passed=%v\n", r.ID, r.Passed)
	}
	if len(report.Cases) == 0 {
		fail(fmt.Errorf("no selected cases"))
	}
	if err := json.NewEncoder(os.Stdout).Encode(report); err != nil {
		fail(err)
	}
	if !allPassed {
		os.Exit(1)
	}
}
