package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// Each DATA represents ONE line. The number of columns for the slice is NCOL
type DATA struct {
	Cols []float64
}

// Parse the file and extract the lines into []DATA
func ParseData(filename, comment string, firstLineIsLegend bool, xcol, ycol int) ([]DATA, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	ncol := -1

	var legend []string = nil
	lines := make([]DATA, 0)
	// Read the lines
	r := bufio.NewReader(file)
	isFirstLine := true
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		// Maybe drop the first line as the titles of columns (legend)
		if isFirstLine {
			isFirstLine = false
			if firstLineIsLegend {
				l := strings.TrimSpace(string(line))
				legend = strings.Fields(l)
				legend[0] = strings.TrimPrefix(legend[0], comment)
				continue
			}
		}
		// Maybe drop comment
		l := strings.TrimSpace(string(line))
		if strings.HasPrefix(l, comment) {
			continue
		}
		fields := strings.Fields(string(line))
		if ncol == -1 {
			ncol = len(fields)
			if xcol >= ncol || ycol >= ncol {
				return nil, nil, errors.New("Not enough colums in file " + filename)
			}
		}

		if len(fields) < ncol {
			return nil, nil, errors.New("Bad formatted line : " + string(line))
		}
		col := make([]float64, ncol)
		for i := range col {
			col[i], err = strconv.ParseFloat(fields[i], 64)
			if err != nil {
				return nil, nil, err
			}
		}
		lines = append(lines, DATA{Cols: col})
	}
	// sort the data according to timestamp
	// sort.SliceStable(lines, func(i, j int) bool {
	// 	return lines[i].Nb < lines[j].Nb
	// })
	// return the results
	return lines, legend, nil
}
