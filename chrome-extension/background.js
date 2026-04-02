const TARGET_SITES = [
  "youtube.com",
  "twitter.com",
  "x.com",
  "instagram.com",
  "threads.net",
  "threads.com"
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

    // Normalize sites
    if (url.includes("youtube.com") || url.includes("youtu.be")) {
      site = "youtube";
    } else if (url.includes("twitter.com") || url.includes("x.com")) {
      site = "twitter";
      if (url.includes("x.com")) url = url.replace("x.com", "twitter.com");
    } else if (url.includes("instagram.com")) {
      site = "instagram";
    } else if (url.includes("threads.net") || url.includes("threads.com")) {
      site = "threads";
    }

    // For Threads and Instagram, try to find the video URL on the page
    if (site === "threads" || site === "instagram") {
      try {
        const results = await chrome.scripting.executeScript({
          target: { tabId: tab.id },
          func: () => {
            const video = document.querySelector('video');
            if (video) {
              return video.src || (video.querySelector('source') ? video.querySelector('source').src : null);
            }
            return null;
          }
        });
        
        if (results && results[0] && results[0].result) {
          url = results[0].result;
          console.log("Found video URL via script:", url);
        }
      } catch (err) {
        console.error("Failed to execute script:", err);
      }
    }

    // If we have a direct video URL from a known site, ensure it's categorized correctly
    if (url.includes("fbcdn.net") && site === "unknown") {
      if (tab.url.includes("threads.net") || tab.url.includes("threads.com")) site = "threads";
      else if (tab.url.includes("instagram.com")) site = "instagram";
    }

    const cookies = await getCookies(url);
    const payload = {
      url: url,
      site: site,
      cookies: formatCookies(cookies),
      userAgent: navigator.userAgent,
      referer: tab.url
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
