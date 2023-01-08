#!/bin/bash -eu

td="$(mktemp -d)"
curl -L https://github.com/gohugoio/hugo/releases/download/v0.109.0/hugo_extended_0.109.0_Linux-64bit.tar.gz | tar xzf - -C "${td}" hugo
chmod 0755 "${td}/hugo"

mkdir -p website/content/projects

for dir in $(find . -maxdepth 1 -type d -not -name '.*' -not -name website -not -name primers -not -name workbooks); do
  cp -r "${dir}" website/content/projects/
done

cp -r primers website/content/

mkdir -p website/content/about
cp CONTRIBUTING.md website/content/about/contributing.md

cp README.md website/content/projects/_index.md

for file in $(find website/content/projects -name README.md); do
  mv "${file}" "${file%README.md}index.md"
done

cd website
"${td}/hugo"
