#!/bin/sh

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for the backend to be up, if we know where it is.
/bin/wait-for-it.sh postgres:5432 -t 60 --

# Run the main container command.
exec "$@"