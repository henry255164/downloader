const TARGET_SITES = [
  "youtube.com",
  "twitter.com",
  "x.com",
  "instagram.com"
];

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
    expirationDate: cookie.expirationDate
  }));
}

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === "download-with-go") {
    const url = info.linkUrl || info.srcUrl || info.pageUrl;
    let site = "unknown";

    if (url.includes("youtube.com") || url.includes("youtu.be")) {
      site = "youtube";
    } else if (url.includes("twitter.com") || url.includes("x.com")) {
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

    console.log("Sending payload to Go:", payload);
    
    // chrome.runtime.sendNativeMessage(
    //   "com.media.downloader.go",
    //   payload,
    //   (response) => {
    //     if (chrome.runtime.lastError) {
    //       console.error("Native Messaging Error:", chrome.runtime.lastError.message);
    //     } else {
    //       console.log("Native Messaging Response:", response);
    //     }
    //   }
    // );
  }
});
