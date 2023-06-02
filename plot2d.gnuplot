#!/usr/bin/gnuplot --persist

set terminal png size 1600,800
set output ARG3

set hidden3d
set isosamples 30

set multiplot layout 1, 2

set label 1
  splot ARG2 title 'groun_truth' lw 3, \
	    ARG1 title 'interpolated' lt 7 lw 1
unset label 1


unset multiplot