document.addEventListener("DOMContentLoaded", () => {
  const currentTheme = localStorage.getItem('theme');
  const themeToggle = document.getElementById('theme-toggle');

  if (currentTheme) {
    document.documentElement.setAttribute('data-theme', currentTheme);
    themeToggle.checked = currentTheme === 'light' ? true : false;
  }

  function toggleTheme(e) {
    const theme = e.target.checked ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }

  themeToggle.addEventListener('change', toggleTheme, false);
});