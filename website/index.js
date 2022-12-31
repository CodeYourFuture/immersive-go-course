// Description: I'm copying all the projects
// It is synchronous as it's only run once when the static site is built.

const fs = require("fs-extra");
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

// filter out everything else that isn't a directory
const directories = availableItems.filter((item) => {
  return fs.statSync(path.join("../", item)).isDirectory();
});

//copy all the directories to the website/content/projects folder
directories.forEach((project) => {
  fs.copySync(
    path.join("../", project),
    path.join(__dirname, "content", "projects", project)
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
    console.log("fetch has failed no contributors to write to file");
  }
};
// I'm using the github api to get the contributors
// using new native fetch api in node
const getContributors = async () => {
  const res = await fetch(
    "https://api.github.com/repos/CodeYourFuture/immersive-go-course/contributors"
  )
    .then((response) => response.json())
    .then((data) => writeToContributors(data))
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
  {
    src: "https://api.github.com/repos/CodeYourFuture/immersive-go-course/contributors",
    target: "contributors.json",
    hugoDir: "data",
  },
];

githubData.forEach((data) =>
  getGithubData(data.src, data.hugoDir, data.subDir, data.target)
);