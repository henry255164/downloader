## 1. Chrome Extension 基礎與選單

- [x] 1.1 初始化 Chrome Extension Manifest V3 結構。
- [x] 1.2 實作 Context Menu for Media Detection，在 YouTube/Twitter/IG 上顯示「Download with Go」。
- [x] 1.3 根據右鍵點擊對象 Extract Media Metadata 並識別 site 類型。

## 2. 會話與 Cookie 管理

- [x] 2.1 實作 Cookie Extraction 機制，擷取目標網域的所有 Cookies。
- [x] 2.2 實作 Secure Cookie Packaging 將 Cookies 格式化為 JSON 準備傳輸。
- [ ] 2.3 整合 Cookies 傳遞機制，將 Metadata 與 Cookies 封裝為完整 Payload。

## 3. Go 本地服務開發

- [ ] 3.1 建立 Go 專案結構並實作 Native Messaging Establishment 與 Message Protocol Compliance (4-byte header)。
- [ ] 3.2 採用 Native Messaging 通訊協議處理 STDIN/STDOUT。
- [ ] 3.3 使用 Golang 實作並行下載 Worker Pool Concurrency 併發管理邏輯。
- [ ] 3.4 實作 Target Path Classification 以及動態路徑映射與環境變數配置。
- [ ] 3.5 實作 yt-dlp Execution 下載核心調用，傳遞 URL 與 Cookie 檔案。

## 4. 系統整合與配置

- [ ] 4.1 撰寫 .env 配置範例，包含 BASE_DOWNLOAD_PATH 與並行限制。
- [ ] 4.2 提供 Native Messaging Manifest 安裝腳本以註冊 Go 服務。
- [ ] 4.3 進行全流程測試：從 Extension 右鍵點擊到 Go 服務完成下載。
