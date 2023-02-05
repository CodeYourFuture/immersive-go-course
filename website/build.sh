#!/bin/bash -eu

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

find website/content -name '*.md' -print0 | xargs -0 "${sed_i[@]}" -e '/^<!--forhugo$/d' -e '/^forhugo-->$/d'

cd website
"${hugo}"
