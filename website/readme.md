## CYF+ Immersive Engineering Programme website

This is a super small static Hugo site. Assets are handled by Hugo Pipes and there are no node modules. 5kb CSS, 500 bytes of JS (brotli'd). Please be a good citizen of the repo and keep this website simple and small.

## Important, parent repo is source of truth

Content pages are readmes from the root folder copied in to Hugo on build. If you want to change the content, open a PR there.

### first run

You need Hugo extended edition, so if you don't have it, install it
https://gohugo.io/installation/

```zsh
brew install hugo
```

### to populate the content folder

```zsh
./website/build.sh
```

### to develop

```zsh
hugo serve
```

### to stage

Open a PR to main and Netlify will create a branch preview for you.

### to deploy

Site deploys automatically when you merge your PR to main.
