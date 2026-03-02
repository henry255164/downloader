## Context

目前使用者需要手動複製網址並使用外部工具下載媒體，過程繁瑣且常受限於登入牆（如 Instagram）或影音分離技術（如 YouTube）。本設計旨在建立一個無縫的 Chrome Extension 與 Go 本地服務聯動系統。

## Goals / Non-Goals

**Goals:**
- 提供右鍵選單直接觸發下載。
- 自動帶入瀏覽器 Cookies 解決登入牆問題。
- 透過 Go 語言實作高效率、具備並行限制（Worker Pool）的下載管理。
- 支援不同網站自動分類儲存。

**Non-Goals:**
- 不在瀏覽器內實作 FFmpeg 合併（由 `yt-dlp` 在本地端處理）。
- 本階段不實作下載進度條回傳至 Extension UI。
- 不提供網頁內嵌的下載按鈕（僅使用右鍵選單）。

## Decisions

### 採用 Native Messaging 通訊協議
**Rationale:** 相比於 Local Web Server，Native Messaging 是 Chrome 官方推薦與本地應用通訊的方式。它更安全（不需要開放埠口）、啟動更自動化（Chrome 會自動啟動 Go 服務）。
**Alternatives:** Local HTTP Server (FastAPI/Go), Clipboard Listener.

### 使用 Golang 實作並行下載 Worker Pool
**Rationale:** Go 的 Goroutines 與 Channels 適合實作輕量級的併發管理。我們將建立一個固定數量的 Worker 群組，從任務 Channel 中領取任務。
**Alternatives:** Python `asyncio`, Rust `tokio`.

### 整合 Cookies 傳遞機制
**Rationale:** Extension 透過 `chrome.cookies.getAll` 抓取內容，轉換為 JSON 格式傳給 Go。Go 將其寫入暫存檔案並在 `yt-dlp` 執行時透過 `--cookies` 參數引用。
**Alternatives:** `--cookies-from-browser` (有檔案權限鎖定風險)。

### 動態路徑映射與環境變數配置
**Rationale:** 使用 `.env` 設定 `BASE_DOWNLOAD_PATH`。Go 服務接收到任務時，根據 `site` 欄位（youtube/twitter/instagram）自動在該路徑下建立子資料夾。

## Risks / Trade-offs

- **[Risk] yt-dlp 缺失** → **Mitigation**: 在 Go 服務啟動時檢查 `yt-dlp` 是否在 PATH 中，若無則報錯。
- **[Risk] Cookie 格式不相容** → **Mitigation**: Go 端需要一個轉換器，將 Chrome 的 JSON Cookie 轉為 `yt-dlp` 可接受的 Netscape 格式或暫存為標準 JSON。
- **[Risk] Native Messaging 安裝複雜** → **Mitigation**: 提供一個安裝腳本（bash/ps1）自動註冊 Manifest 檔案。
