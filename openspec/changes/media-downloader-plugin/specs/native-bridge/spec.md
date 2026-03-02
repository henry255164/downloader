## ADDED Requirements

### Requirement: Native Messaging Establishment
The Chrome Extension SHALL establish a persistent communication channel with the local Go Host Service using the `chrome.runtime.connectNative` API.

#### Scenario: Extension starting up
- **WHEN** the browser starts or the Extension is loaded
- **THEN** the Extension SHALL attempt to connect to the native host `com.media.downloader.go`

### Requirement: Message Protocol Compliance
The Go Host Service and the Extension SHALL use a message format with a 4-byte little-endian length prefix followed by a JSON payload.

#### Scenario: Sending download task
- **WHEN** the Extension sends a message to the native host
- **THEN** the first 4 bytes SHALL indicate the length of the JSON message that follows
