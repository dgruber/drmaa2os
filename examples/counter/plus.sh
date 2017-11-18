#!/bin/sh
read line
echo $line+1  | bc 2>&1
