## ADDED Requirements

### Requirement: Installation Guide for Chrome Extension
The project documentation SHALL provide a step-by-step guide on how to install the Chrome Extension in developer mode.

#### Scenario: User installing the extension
- **WHEN** user reads the README.md
- **THEN** they MUST find instructions to load the extension via `chrome://extensions`

### Requirement: Go Host Setup Instructions
The documentation SHALL specify how to build and register the Go Host service using the provided install script.

#### Scenario: Setting up the native host
- **WHEN** user wants to connect the extension to the local service
- **THEN** the documentation SHALL instruct them to run `./install.sh <EXTENSION_ID>`

### Requirement: Configuration Reference
The documentation SHALL list all available environment variables and their purposes in a clear reference section.

#### Scenario: Configuring the downloader
- **WHEN** user opens the README.md
- **THEN** they SHALL see a table or list explaining `BASE_DOWNLOAD_PATH` and `MAX_CONCURRENT_DOWNLOADS`
