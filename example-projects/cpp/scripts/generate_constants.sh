#!/bin/bash

set -euo pipefail

mkdir -p gen

cat >gen/constants.h <<EOF
#include <string>
#include <vector>

namespace constants {
  static std::vector<std::string> names = {
EOF

for name in "$@"; do
  echo "    \"${name}\"," >> gen/constants.h
done

cat >>gen/constants.h <<EOF
  };
}
EOF
