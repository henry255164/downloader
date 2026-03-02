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

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === "download-with-go") {
    let url = info.linkUrl || info.srcUrl || info.pageUrl;
    let site = "unknown";

    if (url.includes("youtube.com") || url.includes("youtu.be")) {
      site = "youtube";
    } else if (url.includes("twitter.com") || url.includes("x.com")) {
      site = "twitter";
    } else if (url.includes("instagram.com")) {
      site = "instagram";
    }

    const cookies = await getCookies(url);
    console.log(`Extracted ${cookies.length} cookies for ${site}`);
    // Task 2.2 will format cookies
  }
});
