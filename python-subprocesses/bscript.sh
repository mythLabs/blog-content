#!/bin/bash

echo "running for $1 times"

x=1
while [ $x -le $1 ]
do
  echo "."
  x=$(( $x + 1 ))
  sleep 1
done