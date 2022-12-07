#!/bin/sh
# Randomly fail a job
if [ $((RANDOM % 2)) -eq 0 ]; then
    echo "Job failed"
    exit 1
else
    echo "Job succeeded"
    exit 0
fi