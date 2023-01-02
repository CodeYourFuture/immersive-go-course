// Description: I'm copying all the projects
// It is synchronous as it's only run once when the static site is built.

const fs = require("fs");
const path = require("path");

// get a list of all items in the root directory
const items = fs.readdirSync("../");

// filter out hidden files and folders we don't want, like this hugo site as that is recursive madness
// not very robust, but it works for now
const availableItems = items.filter((item) => {
  return (
    !item.startsWith(".") &&
    item !== "node_modules" &&
    item !== "website" &&
    item !== "readme-assets"
  );
});

// filter out everything else that isn't a directory, double check!
const directories = availableItems.filter((item) => {
  return fs.statSync(path.join("../", item)).isDirectory();
});

console.log(directories + " copied to content/projects");

// for each directory, copy the readme.md file to the content/projects folder
directories.forEach((dir) => {
  fs.copyFileSync(
    path.join("../", dir, "readme.md"),
    path.join(__dirname, "content", "projects", `${dir}.md`)
  );
});

// grab everything else from github that we don't want to update manually
// using new native fetch in node 18 - https://nodejs.org/api/fetch.html
const getGithubData = async (src, hugoDir, subDir = "", targetFile) => {
  const res = await fetch(src)
    .then((response) => response.text())
    .then((data) => {
      fs.writeFileSync(path.join(__dirname, hugoDir, subDir, targetFile), data);
    })
    .catch((error) => console.log(error));
  return res;
};

const githubData = [
  {
    src: "https://raw.githubusercontent.com/CodeYourFuture/immersive-go-course/master/README.md",
    target: "_index.md",
    hugoDir: "content",
    subDir: "projects",
  },
  {
    src: "https://raw.githubusercontent.com/CodeYourFuture/immersive-go-course/master/CONTRIBUTING.md",
    target: "contributing.md",
    hugoDir: "content",
    subDir: "about",
  },
];

githubData.forEach((data) =>
  getGithubData(data.src, data.hugoDir, data.subDir, data.target)
);
