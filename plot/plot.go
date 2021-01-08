package main

/**
Program used to draw one or several files in the same graphic
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

	"gonum.org/v1/plot/vg"
)

// Print some values while drawing
var PRINT = false

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
//    file1, file2 ...
func main() {
	flag.Usage = Usage
	p := flag.Bool("p", PRINT, "Print the data while drawing")
	nolegend := flag.Bool("nolegend", false, "False if first line is not interpreted as the column titles for the graphic legend")
	output := flag.String("output", "", "Name of the output graphic")
	root := flag.String("root", "", "Root folder of the files to process")
	comment := flag.String("comment", "#", "Comment line prefixe")
	xlabel := flag.String("xlabel", "", "X axis label")
	ylabel := flag.String("ylabel", "", "Y axis label")
	title := flag.String("title", "", "Graphic title")
	xcol := flag.Int("xcol", 0, "X column number")
	ycol := flag.Int("ycol", -1, "Y column number")

	flag.Parse()

	files := checkOptions(flag.Args(), *root, *p, *xcol, *ycol)

	datas := make([][]DATA, len(files))
	legends := make([][]string, len(files))
	for i, file := range files {
		data, legend, err := ParseData(file, *comment, !*nolegend, *xcol, *ycol)
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

	xlab := buildXLabel(*xlabel, "x", legends, *xcol)
	ylab := buildYLabel(*ylabel, "y", legends, *xcol, *ycol)
	titl := builtTitle(*title, files)
	outp := buildOutput(*output, files)

	err := Draw(datas, outp, titl, xlab, ylab, *xcol, *ycol, legends)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Define the graphic file name
func buildOutput(name string, files []string) string {
	if name != "" {
		return name
	}
	base := filepath.Base(files[0])
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext) + ".png"
}

// Define the title
func builtTitle(title string, files []string) string {
	if title != "" {
		return title
	}
	return filepath.Base(files[0])
}

// Define the Y axis label
func buildYLabel(ylabel, defaut string, legends [][]string, xcol, ycol int) string {
	if ylabel != "" {
		return ylabel
	}
	if legends[0] != nil {
		if ycol == -1 {
			if xcol == 0 {
				label := legends[0][1]
				for i := 2; i < len(legends[0]); i++ {
					label = label + "-" + legends[0][i]
				}
				return label
			} else {
				label := legends[0][0]
				for i := range legends[0] {
					if i == xcol {
						continue
					}
					label = label + "-" + legends[0][i]
				}
				return label
			}
		} else {
			return legends[0][ycol]
		}
	}
	return defaut
}

// Define the X axis label
func buildXLabel(xlabel, defaut string, legends [][]string, xcol int) string {
	if xlabel != "" {
		return xlabel
	}
	if legends[0] != nil {
		return legends[0][xcol]
	}
	return defaut
}

// Draw the lines on the same graphic
func Draw(datas [][]DATA, png, title, xlabel, ylabel string, xcol, ycol int, legends [][]string) error {
	p, err := NewPlot(title, xlabel, ylabel)
	if err != nil {
		return err
	}
	leg := ""
	color := 0
	for j, d := range datas {
		x := make([]float64, len(d))
		y := make([]float64, len(d))
		for i := range d {
			x[i] = d[i].Cols[xcol]
		}
		if ycol == -1 { // Draw all Y columns
			for c := range d[0].Cols {
				if c == xcol {
					continue // Don't draw xcol == ycol
				}
				if legends[j] != nil {
					leg = legends[j][c]
				}
				for i := range d {
					y[i] = d[i].Cols[c]
				}
				Print(x, y, title+" "+leg)
				if err = AddWithLineXY(x, y, leg, color, p); err != nil {
					return err
				}
				color++
			}
		} else { // Draw only the specified Y column
			if legends[0] != nil {
				leg = legends[0][ycol]
			}
			for i := range d {
				y[i] = d[i].Cols[ycol]
			}
			Print(x, y, title+" "+leg)
			if err = AddWithLineXY(x, y, leg, color, p); err != nil {
				return err
			}
			color++
		}
	}
	// Save the plot to a PNG file.
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, png)
}

// Check the program arguments (options) and exit in case of error
func checkOptions(files []string, root string, p bool, xcol, ycol int) []string {
	if xcol < 0 {
		fmt.Println("Error : x column must be positive : ", xcol)
		os.Exit(1)
	}
	if ycol < -1 {
		fmt.Println("Error : y column must be positive : ", ycol)
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
	PRINT = p
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
