# Changelog
All notable changes to this project will be documented in this file.

## [v2.1.0]

### Added

- Added dark/light theme toggle with glassmorphism design
- Added "Copy All" button to copy all configuration commands at once
- Added Dex favicon for better browser tab identification
- Added descriptive text for each command block explaining its purpose
- Added numbered headings and descriptions for better command organization

### Changed

- Improved dark theme visibility and contrast:
  - Fixed text color for "Select which cluster you require a token for:" in dark theme
  - Fixed hover states for tabs and buttons to ensure readable text in dark theme
  - Improved copy button visibility in dark theme with better contrast
  - Fixed "Login Again" button text color to match theme (was blue in dark theme)
- Redesigned UI with glassmorphism effects:
  - Removed top navbar for cleaner interface
  - Styled "Login Again" button with glassmorphism design
  - Improved button and tab hover effects
- Improved command blocks:
  - Separated cluster configuration commands into individual fields
  - Reduced font size for commands and tokens (12px) for better screen fit
  - Moved copy buttons inside command blocks instead of beside them
  - Added spacing and visual separation between command blocks
- Added the ability to pass the cluster certificate directly in the kubeconfig. Lens on macOS runs in its own environment and doesn’t have access to the certificate file.

### Fixed

- Fixed text readability issues in dark theme
- Fixed copy button positioning and visibility
- Fixed color inconsistencies between light and dark themes
