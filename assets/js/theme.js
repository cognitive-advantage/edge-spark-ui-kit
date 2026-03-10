(function () {
  var storageKey = "ui-theme";
  var root = document.documentElement;

  function normalizeTheme(raw) {
    if (raw === "dark") {
      return "dark";
    }
    return "light";
  }

  function applyTheme(theme) {
    root.setAttribute("data-theme", theme);
    try {
      localStorage.setItem(storageKey, theme);
    } catch (_err) {
      // Ignore storage failures in restricted browsing modes.
    }
  }

  function loadTheme() {
    try {
      var saved = localStorage.getItem(storageKey);
      if (saved) {
        applyTheme(normalizeTheme(saved));
        return;
      }
    } catch (_err) {
      // Ignore storage failures in restricted browsing modes.
    }

    var prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
    applyTheme(prefersDark ? "dark" : "light");
  }

  window.uiKitTheme = {
    set: function (theme) {
      applyTheme(normalizeTheme(theme));
    },
    toggle: function () {
      var current = root.getAttribute("data-theme") || "light";
      applyTheme(current === "dark" ? "light" : "dark");
    }
  };

  loadTheme();
})();
