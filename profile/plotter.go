package profile

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	log "github.com/sirupsen/logrus"
)

// CSVPlotter plots every attribute as VM number increases
func CSVPlotter(inFile, outPath string) {
	var (
		records = readResultCSV(inFile)
		rows    = len(records)
		cols    = len(records[0])
	)

	for col := 0; col < cols; col++ {
		// create a new plot for a metric
		p, err := plot.New()
		if err != nil {
			log.Fatalf("Failed creating plot: %v", err)
		}

		p.X.Label.Text = "VM number"
		p.Y.Label.Text = records[0][col]

		// setup data
		pts := make(plotter.XYs, rows-1)
		vmNum := 4
		for row := 1; row < rows; row++ {
			pts[row-1].X = float64(vmNum)
			value, err := strconv.ParseFloat(records[row][col], 64)
			if err != nil {
				log.Fatalf("Failed parsing string to float: %v", err)
			}
			pts[row-1].Y = value
			vmNum += 4
		}

		err = plotutil.AddLinePoints(p, pts)
		if err != nil {
			log.Fatalf("Failed plotting data: %v", err)
		}

		fileName := filepath.Join(outPath, p.Y.Label.Text+".png")
		if err := p.Save(4*vg.Inch, 4*vg.Inch, fileName); err != nil {
			log.Fatalf("Failed saving plot: %v", err)
		}
	}
}

// retrieve data from csv file
func readResultCSV(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Failed reading file %s: %v", filePath, err)
	}

	return records
}