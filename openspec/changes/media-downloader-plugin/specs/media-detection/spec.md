## ADDED Requirements

### Requirement: Context Menu for Media Detection
The Chrome Extension SHALL register a context menu item that appears when the user right-clicks on a link or a video element on supported platforms (YouTube, Twitter, Instagram).

#### Scenario: Right-click on YouTube video
- **WHEN** user right-clicks on a YouTube video link or the video player
- **THEN** a "Download with Go" option MUST appear in the context menu

### Requirement: Extract Media Metadata
The Extension SHALL extract the target URL and determine the source platform (site type) from the context of the right-click event.

#### Scenario: Extracting Twitter video URL
- **WHEN** user selects "Download with Go" on a Twitter video post
- **THEN** the system SHALL capture the tweet URL and identify the platform as "twitter"
