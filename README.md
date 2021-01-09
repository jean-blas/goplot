# plot

Plot one or several files in .png graphics automatically.

_plot [options] file1.res file2.res file3.res_

_plot [options] *.res_

Enter _plot -h_ for common options usage and defaults.

When a file contains several colums and the _ycol_ option is not defined, then all columns are drawn on the same graphic.

## Automation mode

When several files are entered in command line (ex: _plot file1 file2 file3_ or _plot *.res_):

* if option _-automation_ is not set then all the files are plotted on the same graphic.
* If option _-automation_ is set, then each file is plotted into its own graphic file, with its name extension changed to .png

## Options

The simplest use of _plot_ program is

_plot file.res_

which plots the _file.res_ into a _file.png_ graphic with default values for all options, which are:

* automation : True implies one graphic per file
* comment : the comment prefixe value. Lines beginning with this prefixe are ignored (default #)
* nolegend : no automatic legend drawn. By default, the first line should refer to the column title, and is used to compute the graphic legend.
* output : name of the graphic file. By default, it is the name of the first file, with extension .png (ex: file.res => file.png)
* p : print some data while drawing
* pt : draw points instead of lines
* root : define a root folder to search for files (default .)
* title : graphic title
* xcol : the column numero to use as absissa (default 0)
* xlabel : X axis label
* xlength : X axis length in cm
* ycol : the column numero to use as ordinates (default 1)
* ylabel : Y axis label
* ylength : Y axis length in cm
* ytopleg : Y position of the legend (default to bottom)
