#!/usr/bin/gnuplot --persist

data = ARG1
true = ARG2
out = ARG3

set terminal png size 1200,800
set output out

plot data with lp title 'Sol' lw 4, \
     true with p title 'True' lw 6
