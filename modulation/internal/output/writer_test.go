package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
)

func TestWriterWritesCSV(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "orbit.csv")
	writer := NewWriter(path)
	writer.Write(7, []physics.Body{
		{
			Mass:     1,
			Position: physics.Vector{X: 1, Y: 2, Z: 3},
			Velocity: physics.Vector{X: 4, Y: 0, Z: 0},
		},
	})
	writer.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	got := string(data)
	if !strings.Contains(got, "step,body,x,y,z,energy") {
		t.Fatalf("missing header in %q", got)
	}
	if !strings.Contains(got, "7,0,1.000000e+00,2.000000e+00,3.000000e+00") {
		t.Fatalf("missing body row in %q", got)
	}
}
