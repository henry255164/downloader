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

function normalizeDownloadUrl(url) {
  const mediaUrl = new URL(url);
  if (/\.(?:mp4|webm|mov|m4v)$/i.test(mediaUrl.pathname) &&
      (mediaUrl.searchParams.has('bytestart') || mediaUrl.searchParams.has('byteend'))) {
    mediaUrl.searchParams.delete('bytestart');
    mediaUrl.searchParams.delete('byteend');
  }
  return mediaUrl.href;
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
          func: async () => {
            for (let attempt = 0; attempt < 10; attempt++) {
              for (const video of document.querySelectorAll('video')) {
                const videoUrl = [video.currentSrc, video.src, video.querySelector('source')?.src]
                  .find(src => /^https?:\/\//.test(src));
                if (videoUrl) return videoUrl;
              }
              await new Promise(resolve => setTimeout(resolve, 500));
            }
            const mediaEntry = performance.getEntriesByType('resource').find(entry => {
              try {
                const mediaUrl = new URL(entry.name);
                const host = mediaUrl.hostname;
                return mediaUrl.protocol === 'https:' && mediaUrl.pathname.endsWith('.mp4') &&
                  (host === 'fbcdn.net' || host.endsWith('.fbcdn.net') ||
                   host === 'cdninstagram.com' || host.endsWith('.cdninstagram.com'));
              } catch {
                return false;
              }
            });
            return mediaEntry?.name || null;
          }
        });
        
        if (!results?.[0]?.result && site === "threads") {
          console.error("No direct video URL found; download skipped:", url);
          return;
        }
        if (results?.[0]?.result) {
          url = results[0].result;
          console.log("Found direct video URL via script:", url);
        }
      } catch (err) {
        console.error("Failed to execute script:", err);
        if (site === "threads") return;
      }
    }

    // If we have a direct video URL from a known site, ensure it's categorized correctly
    if (url.includes("fbcdn.net") && site === "unknown") {
      if (tab.url.includes("threads.net") || tab.url.includes("threads.com")) site = "threads";
      else if (tab.url.includes("instagram.com")) site = "instagram";
    }

    url = normalizeDownloadUrl(url);
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
