// Theme Toggle Functionality
(function() {
  // Get theme from localStorage or default to light
  const getTheme = () => {
    const savedTheme = localStorage.getItem('theme');
    return savedTheme || 'light';
  };

  // Apply theme
  const applyTheme = (theme) => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
    updateToggleButton(theme);
  };

  // Update toggle button text
  const updateToggleButton = (theme) => {
    const button = document.getElementById('theme-toggle');
    if (button) {
      button.textContent = theme === 'dark' ? '☀️ Light' : '🌙 Dark';
    }
  };

  // Initialize theme on page load
  document.addEventListener('DOMContentLoaded', () => {
    const theme = getTheme();
    applyTheme(theme);

    // Create toggle button if it doesn't exist
    if (!document.getElementById('theme-toggle')) {
      const button = document.createElement('button');
      button.id = 'theme-toggle';
      button.className = 'theme-toggle';
      button.textContent = theme === 'dark' ? '☀️ Light' : '🌙 Dark';
      button.onclick = () => {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        applyTheme(newTheme);
      };
      document.body.appendChild(button);
    }
  });
})();
