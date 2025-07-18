#!/bin/sh

go test -fuzzminimizetime 0 -fuzz=FuzzMultipartParser

