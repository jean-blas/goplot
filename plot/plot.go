package main

/**
Program used to draw one or several files in one or several graphics
Structure of the files:
	* data are in columns
	* Ignore comment lines (begins with #)
Comparison of all columns in the same file
*/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

var (
	// Print some values while drawing
	PRINT = false
	// Draw points instead of lines
	POINT = false
	// Graphic title
	TITLE = ""
	// Graphic X axis label
	XLABEL = ""
	// Graphic Y axis label
	YLABEL = ""
	// Draw the legend
	NOLEGEND = false
	// Graphic file name
	OUTPUT = ""
	// Legend Y position on top/bottom
	YTOPLEGEND = false
	// Graphic X axis length in cm
	XLENGTH = 10
	// Graphic Y axis length in cm
	YLENGTH = 10
)

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s: utility program used to plot a file automatically\n\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), "plot [options] file1.res file2.res file3.res")
	fmt.Fprintln(flag.CommandLine.Output(), "plot [options] file*")
	fmt.Fprintln(flag.CommandLine.Output(), "plot [options] *.res")
	fmt.Fprintf(flag.CommandLine.Output(), "\nUsage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

// go run main.go parser.go [-output graphic.png] [-root root_folder] [-p]
//    [-xcol 0] [-ycol 1] [-comment #] [-nolegend]
//    [-xlabel nb] [-ylabel size] [-title throughput]
//    [-pt] [-xlength 10] [-ylength 10] [-automation] [-ytopleg]
//    file1, file2 ...
func main() {
	flag.Usage = Usage
	automation := flag.Bool("automation", false, "True implies one graphic per file")
	comment := flag.String("comment", "#", "Comment line prefixe")
	nolegend := flag.Bool("nolegend", NOLEGEND, "False if first line is not interpreted as the column titles for the graphic legend")
	output := flag.String("output", OUTPUT, "Name of the output graphic")
	p := flag.Bool("p", PRINT, "Print the data while drawing")
	pt := flag.Bool("pt", POINT, "Draw points instead of lines")
	root := flag.String("root", "", "Root folder of the files to process")
	title := flag.String("title", TITLE, "Graphic title")
	xcol := flag.Int("xcol", 0, "X column number")
	xlabel := flag.String("xlabel", XLABEL, "X axis label")
	xlength := flag.Int("xlength", XLENGTH, "Graphic X axis length in cm")
	ycol := flag.Int("ycol", -1, "Y column number")
	ylength := flag.Int("ylength", YLENGTH, "Graphic Y axis length in cm")
	ylabel := flag.String("ylabel", YLABEL, "Y axis label")
	ytopleg := flag.Bool("ytopleg", YTOPLEGEND, "Y position of the legend (default to bottom)")

	flag.Parse()

	NOLEGEND = !*nolegend
	PRINT = *p
	POINT = *pt
	YTOPLEGEND = *ytopleg
	XLENGTH = *xlength
	YLENGTH = *ylength

	files := checkOptions(flag.Args(), *root, *xcol, *ycol, *xlength, *ylength)

	datas := make([][]DATA, len(files))
	legends := make([][]string, len(files))
	for i, file := range files {
		data, legend, err := ParseData(file, *comment, NOLEGEND, *xcol, *ycol)
		if err != nil {
			fmt.Println(file, err)
			os.Exit(1)
		}
		datas[i] = data
		legends[i] = legend
		if PRINT {
			fmt.Println(file, "data: ", datas[i])
			fmt.Println(file, "legend: ", legends[i])
		}
	}

	if *automation {
		var wg sync.WaitGroup
		wg.Add(len(files))
		for i := range files {
			go func(i int) {
				xlab := buildXLabel(*xlabel, "x", legends[i], *xcol)
				ylab := buildYLabel(*ylabel, "y", legends[i], *xcol, *ycol)
				titl := builtTitle(*title, files[i])
				outp := buildOutput(*output, files[i])
				if err := drawDataAutomation(datas[i], *xcol, *ycol, legends[i], xlab, ylab, titl, outp); err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	} else {
		XLABEL = buildXLabel(*xlabel, "x", legends[0], *xcol)
		YLABEL = buildYLabel(*ylabel, "y", legends[0], *xcol, *ycol)
		TITLE = builtTitle(*title, files[0])
		OUTPUT = buildOutput(*output, files[0])

		err := drawData(datas, *xcol, *ycol, legends)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// Define the graphic file name
func buildOutput(name string, files string) string {
	if name != "" {
		return name
	}
	base := filepath.Base(files)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext) + ".png"
}

// Define the title
func builtTitle(title string, files string) string {
	if title != "" {
		return title
	}
	return filepath.Base(files)
}

// Define the Y axis label
func buildYLabel(ylabel, defaut string, legends []string, xcol, ycol int) string {
	if ylabel != "" {
		return ylabel
	}
	if legends != nil {
		if ycol == -1 {
			if xcol == 0 {
				label := legends[1]
				for i := 2; i < len(legends); i++ {
					label = label + "-" + legends[i]
				}
				return label
			} else {
				label := legends[0]
				for i := range legends {
					if i == xcol {
						continue
					}
					label = label + "-" + legends[i]
				}
				return label
			}
		} else {
			return legends[ycol]
		}
	}
	return defaut
}

// Define the X axis label
func buildXLabel(xlabel, defaut string, legends []string, xcol int) string {
	if xlabel != "" {
		return xlabel
	}
	if legends != nil {
		return legends[xcol]
	}
	return defaut
}

// draw the lines on several graphics
func drawDataAutomation(datas []DATA, xcol, ycol int, legends []string, xlab, ylab, titl, outp string) error {
	p, err := NewPlot(titl, xlab, ylab)
	if err != nil {
		return err
	}
	x := make([]float64, len(datas))
	y := make([]float64, len(datas))
	for i := range datas {
		x[i] = datas[i].Cols[xcol]
	}
	color := 0
	if ycol == -1 { // Draw all Y columns
		for col := range datas[0].Cols {
			if col == xcol {
				continue // Don't draw xcol == ycol
			}
			innerDraw(datas, x, y, getLegend(legends, col), color, col, p)
			color++
		}
	} else { // Draw only the specified Y column
		innerDraw(datas, x, y, getLegend(legends, ycol), color, ycol, p)
		color++
	}
	// Save the plot to a PNG file.
	err = p.Save(vg.Length(XLENGTH)*vg.Centimeter, vg.Length(YLENGTH)*vg.Centimeter, outp)
	if err != nil {
		return err
	}
	return nil
}

// draw the lines on the same graphic
func drawData(datas [][]DATA, xcol, ycol int, legends [][]string) error {
	p, err := NewPlot(TITLE, XLABEL, YLABEL)
	if err != nil {
		return err
	}
	color := 0
	for j, d := range datas {
		x := make([]float64, len(d))
		y := make([]float64, len(d))
		for i := range d {
			x[i] = d[i].Cols[xcol]
		}
		if ycol == -1 { // Draw all Y columns
			for col := range d[0].Cols {
				if col == xcol {
					continue // Don't draw xcol == ycol
				}
				innerDraw(d, x, y, getLegend(legends[j], col), color, col, p)
				color++
			}
		} else { // Draw only the specified Y column
			innerDraw(d, x, y, getLegend(legends[0], ycol), color, ycol, p)
			color++
		}
	}
	// Save the plot to a PNG file.
	return p.Save(vg.Length(XLENGTH)*vg.Centimeter, vg.Length(YLENGTH)*vg.Centimeter, OUTPUT)
}

func getLegend(legend []string, c int) string {
	if legend != nil {
		return legend[c]
	}
	return ""
}

func innerDraw(d []DATA, x, y []float64, leg string, color, c int, p *plot.Plot) error {
	for i := range d {
		y[i] = d[i].Cols[c]
	}
	Print(x, y, TITLE+" "+leg)
	if POINT {
		if err := AddWithPointsXY(x, y, leg, color, p); err != nil {
			return err
		}
	} else {
		if err := AddWithLineXY(x, y, leg, color, p); err != nil {
			return err
		}
	}
	color++
	return nil
}

// Check the program arguments (options) and exit in case of error
func checkOptions(files []string, root string, xcol, ycol, xlength, ylength int) []string {
	if xcol < 0 {
		fmt.Println("Error : x column must be positive : ", xcol)
		os.Exit(1)
	}
	if ycol < -1 {
		fmt.Println("Error : y column must be positive : ", ycol)
		os.Exit(1)
	}
	if xlength < 0 {
		fmt.Println("Error : x axis length must be positive : ", xlength)
		os.Exit(1)
	}
	if ylength < 0 {
		fmt.Println("Error : y axis length must be positive : ", ylength)
		os.Exit(1)
	}
	if root != "" {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			fmt.Println("Error : folder does not exist : ", root)
			os.Exit(1)
		}
		info, _ := os.Stat(root)
		if !info.IsDir() {
			fmt.Println("Error : this is not a folder : ", root)
			os.Exit(1)
		}
	}

	if len(files) == 0 {
		fmt.Println("Error : no file to process. Exiting...")
		os.Exit(1)
	}

	ofiles := make([]string, 0)
	// first check if some files contain a unix pattern (ex: *.res)
	unxPatterns := make([]string, 0)
	stdFiles := make([]string, 0)
	for _, f := range files {
		if strings.Contains(f, "*") {
			unxPatterns = append(unxPatterns, f)
		} else {
			stdFiles = append(stdFiles, f)
		}
	}
	// load standard files
	for _, f := range stdFiles {
		if !filepath.IsAbs(filepath.Dir(f)) && root != "" {
			f = filepath.Join(root, f)
		}
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Println("File does not exist : ", f)
			os.Exit(1)
		}
		info, _ := os.Stat(f)
		if !IsOwnerReadable(info) {
			fmt.Println("File is not readable : ", f)
			os.Exit(1)
		}
		ofiles = append(ofiles, f)
	}
	// load files compliant with some eventual unix patterns
	if len(unxPatterns) > 0 {
		for _, p := range unxPatterns {
			pp := p
			if !filepath.IsAbs(p) && root != "" {
				pp = filepath.Join(root, p)
			}
			matches, err := filepath.Glob(pp)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			ofiles = append(ofiles, matches...)
		}
	}
	return ofiles
}

// Print some data in a prettier way
func Print(x, y []float64, legend string) {
	if !PRINT {
		return
	}
	fmt.Println(legend)
	for i := range x {
		fmt.Println("\t", x[i], y[i])
	}
}
