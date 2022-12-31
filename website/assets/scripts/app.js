// I haven't written anything in particular as this is a prototype
// My guess is that we will want to hydrate with React
// so these are just a couple of togglers to make the prototype basically function

// menu toggle
const menuToggles = document.querySelectorAll(".js-menu-toggle");
const menu = document.getElementById("site-menu");
const closeMenu = document.getElementById("close-menu");
const skipLink = document.getElementById("skip-link");
function toggleMenu() {
  menu.classList.toggle("is-active");
  menu.toggleAttribute("hidden");
  if (menu.classList.contains("is-active")) {
    history.pushState(null, null, "#site-menu");
    menu.focus();
  } else {
    skipLink.focus();
  }
}
// listeners - we allow anything to toggle Menu if nominated with the class
menuToggles.forEach((toggle) => toggle.addEventListener("click", toggleMenu));

// if the menu is open listen for escape key to close it
window.addEventListener("keydown", (e) => {
  if (e.key === "Escape" && menu.classList.contains("is-active")) {
    toggleMenu();
  }
  // if meta plus slash key combination is pressed toggle the menu
  if (e.key === "/" && e.metaKey) {
    toggleMenu();
  }
});

// throw in an intersection observer for a little delight
const watched = document.querySelectorAll(".is-watched");
const observer = new IntersectionObserver((blocks) => {
  blocks.forEach((block) => {
    if (block.intersectionRatio > 0) {
      block.target.classList.add("is-visible");
      block.target.classList.add("was-triggered");
    } else {
      block.target.classList.remove("is-visible");
    }
  });
});

watched.forEach((block) => {
  observer.observe(block);
});
/**
 * Generic toggler
 * you can pass in a target of the toggle
 * with data-toggle-target in t html
 */
const toggles = document.querySelectorAll(".js-toggle");
function toggleMe(event) {
  const button = event.currentTarget;
  const toggleTarget = document.getElementById(
    event.currentTarget.dataset.toggleTarget
  );

  if (toggleTarget) {
    toggleTarget.classList.toggle("is-active");
  }
  button.classList.toggle("is-active");
}
// listener
toggles.forEach(function (toggle) {
  toggle.addEventListener("click", toggleMe);
});
