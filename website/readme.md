## CYF+ Immersive Engineering Programme website

This is a super small static Hugo site. Assets are handled by Hugo Pipes and there are no node modules. 5kb CSS, 500 bytes of JS (brotli'd). Please be a good citizen of the repo and keep this website simple and small.

## Important, parent repo is source of truth

Project pages are the readmes from the parent repo copied in to Hugo on build. If you want to change the project page content, open a PR there.

The other content -- /workbooks, /primers, /about -- has been exported from Google Docs and placed in version control here. Let me know if you have a better idea!

### first run

You need Hugo extended edition, so if you don't have it, install it
https://gohugo.io/installation/

```zsh
brew install hugo
```

You need Node ^18 LTS so if you're not running that

```zsh
nvm use 18.12.1
```

### to develop

```zsh
yarn dev
```

That's it. It's probably running on `/localhost:1313`
