# Media Downloader Plugin

這是一個強大的 Chrome 擴充功能，能夠從 YouTube、Twitter (X) 和 Instagram 等多個平台偵測媒體資源（影片與圖片），並透過本地端的 Go 服務利用 `yt-dlp` 進行高效的多工下載。

## 核心功能

- **多平台支援**: 支援 YouTube, Twitter, Instagram。
- **Cookies 整合**: 自動抓取瀏覽器會話 Cookies，解決登入牆與高品質影片下載限制。
- **高效下載**: 本地 Go 服務實作 Worker Pool，支援並行下載任務。
- **自動分類**: 根據來源網站自動將檔案儲存至對應資料夾。

---

## 安裝指南

### 1. Chrome 擴充功能 (Extension)

1.  開啟 Chrome 瀏覽器，進入 `chrome://extensions/`。
2.  開啟右上角的 **「開發者模式」 (Developer mode)**。
3.  點擊 **「載入解壓縮的擴充功能」 (Load unpacked)**。
4.  選擇專案中的 `chrome-extension/` 目錄。
5.  安裝完成後，請**複製擴充功能的「識別碼」 (Extension ID)**（例如：`abcdefghijklmnopqrstuvwxyz`），稍後安裝 Go Host 時會用到。

### 2. 本地 Go 服務 (Host)

確保您的系統已安裝 `go` 語言環境以及 `yt-dlp` 下載工具。

1.  開啟終端機並進入 `downloader-host/` 目錄：
    ```bash
    cd downloader-host
    ```
2.  執行安裝與註冊腳本（需傳入剛才複製的 Extension ID）：
    ```bash
    ./install.sh <YOUR_EXTENSION_ID>
    ```
    *此腳本會編譯 Go 程式並將 Native Messaging Manifest 註冊到 Chrome 的系統路徑中。*

---

## 配置說明 (.env)

在 `downloader-host/` 目錄下建立 `.env` 檔案（可參考 `.env.example`）：

| 設定項 | 說明 | 預設值 |
| :--- | :--- | :--- |
| `BASE_DOWNLOAD_PATH` | 下載檔案儲存的根目錄路徑 | `./downloads` |
| `MAX_CONCURRENT_DOWNLOADS` | 同時進行下載的最大任務數量 (Worker Pool 大小) | `3` |

---

## 使用方式

1.  瀏覽 YouTube、Twitter 或 Instagram。
2.  在影片連結、影片播放器或貼文上點擊 **右鍵**。
3.  選擇選單中的 **「Download with Go」**。
4.  本地 Go 服務會自動接收任務、抓取 Cookies 並開始下載。
5.  檔案將依網站類型存放在 `BASE_DOWNLOAD_PATH` 下的子目錄中（例如：`./downloads/youtube/`）。

---

## 技術架構

- **前端**: Chrome Extension (Manifest V3, Context Menus, Cookies API)。
- **通訊**: Chrome Native Messaging Protocol (4-byte length prefix)。
- **後端**: Golang (Worker Pool, Exec command)。
- **核心**: `yt-dlp`。

## 開發者

由 Gemini CLI 協助開發。
