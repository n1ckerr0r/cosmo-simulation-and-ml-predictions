package output

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/n1ckerr0r/cosmo-simulation-and-ml-predictions/modulation/internal/physics"
)

type Writer struct {
	file   *os.File
	writer *csv.Writer
}

func NewWriter(path string) *Writer {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	writer := csv.NewWriter(file)

	writer.Write([]string{"step", "body", "x", "y", "z", "energy"})

	return &Writer{file: file, writer: writer}
}

func (w *Writer) Write(step int, bodies []physics.Body) {
	energy := physics.TotalEnergy(bodies)

	for i, b := range bodies {
		w.writer.Write([]string{
			strconv.Itoa(step),
			strconv.Itoa(i),
			toStr(b.Position.X),
			toStr(b.Position.Y),
			toStr(b.Position.Z),
			toStr(energy),
		})
	}
}

func (w *Writer) Close() {
	w.writer.Flush()
	w.file.Close()
}

func toStr(v float64) string {
	return strconv.FormatFloat(v, 'e', 6, 64)
}

