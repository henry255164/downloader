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

chrome.contextMenus.onClicked.addListener((info, tab) => {
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

    console.log("Download task:", { url, site });
    // Next tasks will handle cookie extraction and sending to Go
  }
});
