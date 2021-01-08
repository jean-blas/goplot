# plot

Plot one or several files in an .png graphics automatically.

_plot [options] file1.res file2.res file3.res_

_plot [options] *.res_

Enter _plot -h_ for common options usage and defaults.

When a file contains several colums and the _ycol_ option is not defined, then all columns are drawn on the same graphics.

When several files are entered in command line (ex: _plot file1 file2 file3_) then the files are plotted on the same graphics as well.

## Options

The simplest use of _plot_ program is

_plot file.res_

which plots the _file_ in a _file.png_ graphics with default values for all options, which are:

* root : define a root folder to search for files (default .)
* xcol : the column numero to use as absissa (default 0)
* ycol : the column numero to use as ordinates (default 1)
* comment : the comment prefixe value (default #)
* nologend : no automatic legend drawn. By default, the first line should refer to the column title, and is used to compute the graphics legend.
* xlabel : X axis label
* ylabel : Y axis label
* title : graphics title
* p : print some data while drawing
* output : name of the graphics file. By default, it is the name of the first file, with extension .png (ex: file.res => file.png)
