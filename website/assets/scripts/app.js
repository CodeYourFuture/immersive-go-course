// I haven't written anything fancy here, just a few things to make the site work
document.addEventListener("DOMContentLoaded", () => {
  // menu toggle
  const menuToggles = document.querySelectorAll(".js-menu-toggle");
  const menu = document.getElementById("site-menu");
  const skipLink = document.getElementById("skip-link");
  function toggleMenu() {
    menu.classList.toggle("is-active");
    menu.toggleAttribute("hidden");
    if (menu.getAttribute("hidden") == null) {
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

  // dark mode toggle
  const darkModeToggle = document.getElementById("mode-toggle");
  darkModeToggle.addEventListener("click", () => {
    document.body.classList.toggle("is-dark-mode");
    document.body.classList.toggle("is-light-mode");
  });

  //details toggle
  const toggleDetails = document.getElementById("toggle-details");

  if (toggleDetails) {
    toggleDetails.addEventListener("click", () => {
      document.querySelectorAll("details").forEach((detail) => {
        detail.open = !detail.open;
      });
    });
    if (window.location.hash === "#toggle-details") {
      toggleDetails.click(); // send people straight to the toggled details
    }
  }
});
