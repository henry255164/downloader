## Why

提供一個強大的 Chrome 擴充功能，能夠從 YouTube、Twitter (X) 和 Instagram 等多個平台抓取媒體資源（影片與圖片），並透過本地 Go 服務進行高效的多工下載。這解決了瀏覽器對檔案系統限制、大檔案下載不穩定以及特定平台（如 YouTube）影音分離與登入驗證（Cookies）的技術挑戰。

## What Changes

- 新增 Chrome Extension 前端介面，包含右鍵選單觸發下載。
- 實作 Cookies 抓取機制，自動帶入當前瀏覽器的 Session 資訊至下載任務。
- 新增以 Golang 撰寫的 Native Messaging Host，負責接收任務並執行下載。
- 整合 `yt-dlp` 作為核心下載引擎。
- 實作 Worker Pool 併發管理，支援多工下載。
- 支援透過 `.env` 設定下載路徑與併發限制，並根據網站自動分類。

## Capabilities

### New Capabilities

- `media-detection`: 負責在目標網站（YouTube, Twitter, IG）偵測媒體資源並提供右鍵選單觸發。
- `session-management`: 負責從瀏覽器擷取必要的 Cookies 並安全傳遞給下載服務。
- `native-bridge`: 建立 Chrome Extension 與本地 Go 服務之間的通訊橋接。
- `download-service`: 本地 Go 服務，負責 Worker Pool 管理、路徑分配及呼叫 `yt-dlp`。

### Modified Capabilities

(none)

## Impact

- **Affected code**: 
    - Chrome Extension (Manifest V3, Background Service Worker, Content Scripts).
    - Go Host Service (main.go, worker pool logic, native messaging wrapper).
- **Dependencies**: `yt-dlp` (must be installed on host system).
- **Security**: 需要存取 `cookies` 權限以及 `nativeMessaging` 權限。
