#!/bin/bash -eu

SCRIPT_DIR="$(cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}" && cd .. && pwd)"

cd "${REPO_ROOT}"

case "$(uname)" in
  Darwin)
    sed_i=("sed" "-i" "")
    ;;
  Linux)
    sed_i=("sed" "-i")
    ;;
  *)
    echo >&2 "Unrecognised uname: $(uname) - script assumes Linux or Darwin"
    exit 1
    ;;
esac

td="$(mktemp -d)"
case "$(uname)" in
  Darwin)
    if command -v hugo >/dev/null 2>/dev/null; then
      hugo=hugo
    else
      curl -L https://github.com/gohugoio/hugo/releases/download/v0.109.0/hugo_extended_0.109.0_darwin-universal.tar.gz | tar xzf - -C "${td}" hugo
      chmod 0755 "${td}/hugo"
      hugo="${td}/hugo"
    fi
    ;;
  Linux)
    curl -L https://github.com/gohugoio/hugo/releases/download/v0.109.0/hugo_extended_0.109.0_Linux-64bit.tar.gz | tar xzf - -C "${td}" hugo
    chmod 0755 "${td}/hugo"
    hugo="${td}/hugo"
    ;;
esac

cp -r prep primers projects website/content/

cp CONTRIBUTING.md website/content/about/contributing.md

mkdir -p website/data/projects
cp projects/metadata.json website/data/projects/metadata.json

mv website/content/projects/README.md website/content/projects/_index.md

for file in $(find website/content/projects -name README.md); do
  mv "${file}" "${file%README.md}index.md"
done

# Rename README.md based on directory structure within primers so we can have bundles and images stay working
find website/content/primers -type d | while read -r dir; do
  if [[ -f "${dir}/README.md" ]]; then
    # Check if the directory has subdirectories
    if find "${dir}" -mindepth 1 -type d | read; then
      # If it has subdirectories, rename README.md to _index.md
      mv "${dir}/README.md" "${dir}/_index.md"
    else
      # If it does not have subdirectories, rename README.md to index.md
      mv "${dir}/README.md" "${dir}/index.md"
    fi
  fi
done

find website/content -name '*.md' -print0 | xargs -0 "${sed_i[@]}" -e '/^<!--forhugo$/d' -e '/^forhugo-->$/d'

cd website
"${hugo}"
