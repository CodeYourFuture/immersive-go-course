// I haven't written anything fancy here, just a few things to make the site work

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
