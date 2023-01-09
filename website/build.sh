#!/bin/bash -eu

if "$(command -v hugo)" >/dev/null 2>/dev/null; then
  hugo=hugo
else
  td="$(mktemp -d)"
  case "$(uname)" in
    Darwin)
      curl -L https://github.com/gohugoio/hugo/releases/download/v0.109.0/hugo_extended_0.109.0_darwin-universal.tar.gz | tar xzf - -C "${td}" hugo
      ;;
    Linux)
      curl -L https://github.com/gohugoio/hugo/releases/download/v0.109.0/hugo_extended_0.109.0_Linux-64bit.tar.gz | tar xzf - -C "${td}" hugo
      ;;
  esac
  chmod 0755 "${td}/hugo"
  hugo="${td}/hugo"
fi

cp -r prep primers projects website/content/

cp CONTRIBUTING.md website/content/about/contributing.md

cp projects/metadata.json website/data/projects/metadata.json

mv website/content/projects/README.md website/content/projects/_index.md

for file in $(find website/content/projects -name README.md); do
  mv "${file}" "${file%README.md}index.md"
done

cd website
"${hugo}"
