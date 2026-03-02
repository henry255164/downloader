## ADDED Requirements

### Requirement: Cookie Extraction
The Chrome Extension SHALL use the `chrome.cookies` API to extract all cookies for the current domain where the media resource is located.

#### Scenario: Extracting Instagram session cookies
- **WHEN** user initiates a download on Instagram
- **THEN** the system SHALL capture the `sessionid` and other relevant cookies for `instagram.com`

### Requirement: Secure Cookie Packaging
The system SHALL format the extracted cookies into a `yt-dlp` compatible format (e.g., JSON or Netscape) for transmission to the download service.

#### Scenario: Packaging for Go service
- **WHEN** cookies are successfully extracted
- **THEN** they MUST be bundled into a JSON payload with the download URL and site metadata
