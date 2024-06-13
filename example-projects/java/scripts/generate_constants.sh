#!/bin/bash

set -euo pipefail

mkdir -p com/example/gen

cat >com/example/gen/Constants.java <<EOF
package com.example.gen;

import java.util.ArrayList;

public class Constants {
  public static Iterable<String> names() {
    ArrayList<String> names = new ArrayList<>();
EOF

for name in "$@"; do
  echo "    names.add(\"${name}\");" >>com/example/gen/Constants.java
done

cat >>com/example/gen/Constants.java <<EOF
    return names;
  }
}
EOF
