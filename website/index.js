// Description: I'm copying all the projects
// It is synchronous as it's only run once when the static site is built.

const fs = require("fs-extra");
const path = require("path");

// get a list of all items in the root directory
const items = fs.readdirSync("../");

// filter out hidden files and folders we don't want, like this hugo site as that is recursive madness
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
// rename all the README.md files to _index.md
const copiedProjects = fs.readdirSync(
  path.join(__dirname, "content", "projects")
);

copiedProjects.forEach((folder) => {
  if (!folder.startsWith(".")) {
    fs.renameSync(
      path.join(__dirname, "content", "projects", folder, "README.md"),
      path.join(__dirname, "content", "projects", folder, "_index.md")
    );
  }
});
