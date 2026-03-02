const TARGET_SITES = [
  "youtube.com",
  "twitter.com",
  "x.com",
  "instagram.com"
];

let nativePort = null;

// Ensure we have a persistent connection to the Go host
function getNativePort() {
  if (nativePort) return nativePort;

  nativePort = chrome.runtime.connectNative("com.media.downloader.go");
  
  nativePort.onMessage.addListener((response) => {
    console.log("Native Messaging Response:", response);
  });

  nativePort.onDisconnect.addListener(() => {
    console.error("Native Messaging Port Disconnected:", chrome.runtime.lastError?.message);
    nativePort = null;
  });

  return nativePort;
}

chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: "download-with-go",
    title: "Download with Go",
    contexts: ["link", "video", "page"],
    documentUrlPatterns: TARGET_SITES.map(site => `*://*.${site}/*`)
  });
});

async function getCookies(url) {
  const domain = new URL(url).hostname;
  return new Promise((resolve) => {
    chrome.cookies.getAll({ domain }, (cookies) => {
      resolve(cookies);
    });
  });
}

function formatCookies(cookies) {
  return cookies.map(cookie => ({
    name: cookie.name,
    value: cookie.value,
    domain: cookie.domain,
    path: cookie.path,
    secure: cookie.secure,
    httpOnly: cookie.httpOnly,
    expirationDate: cookie.expirationDate,
    hostOnly: cookie.hostOnly
  }));
}

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === "download-with-go") {
    let url = info.linkUrl || info.srcUrl || info.pageUrl;
    let site = "unknown";

    // Normalize x.com to twitter.com for better yt-dlp compatibility
    if (url.includes("x.com")) {
      url = url.replace("x.com", "twitter.com");
    }

    if (url.includes("youtube.com") || url.includes("youtu.be")) {
      site = "youtube";
    } else if (url.includes("twitter.com")) {
      site = "twitter";
    } else if (url.includes("instagram.com")) {
      site = "instagram";
    }

    const cookies = await getCookies(url);
    const payload = {
      url: url,
      site: site,
      cookies: formatCookies(cookies)
    };

    console.log("Sending payload to Go (via Persistent Port):", payload);
    
    try {
      const port = getNativePort();
      port.postMessage(payload);
    } catch (err) {
      console.error("Failed to send message:", err);
    }
  }
});
