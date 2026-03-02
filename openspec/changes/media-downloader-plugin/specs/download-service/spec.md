## ADDED Requirements

### Requirement: Worker Pool Concurrency
The Go Host Service SHALL implement a worker pool with a configurable maximum number of concurrent downloads.

#### Scenario: Submitting multiple tasks
- **WHEN** 5 download tasks are sent while the pool size is 2
- **THEN** exactly 2 SHALL download concurrently while 3 remain in the queue

### Requirement: Target Path Classification
The Go Host Service SHALL organize downloaded files into platform-specific directories (e.g., `youtube/`, `instagram/`) within a base directory defined in `.env`.

#### Scenario: Downloading from Instagram
- **WHEN** an Instagram media link is processed
- **THEN** the file SHALL be saved in the `instagram/` subdirectory of the base path

### Requirement: yt-dlp Execution
The Go Host Service SHALL invoke `yt-dlp` for each task, passing the URL and the extracted cookies as arguments.

#### Scenario: Successful YouTube download
- **WHEN** a YouTube task reaches the top of the queue
- **THEN** the system SHALL execute `yt-dlp --cookies <temp_cookie_path> <url>` to perform the download
