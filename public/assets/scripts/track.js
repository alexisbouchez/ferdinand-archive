((window) => {
  const {
    screen: { width, height },
    navigator: { language },
    location,
    localStorage,
    document,
    history,
  } = window;
  const { hostname, href } = location;
  const { currentScript, referrer } = document;
  console.log("window", window);
  let initialized = false;
  console.log("href: " + href);
  console.log("hostname: " + hostname);
  console.log("location", location);
  console.log("data-website-id", currentScript.dataset.websiteId);
  console.log("referrer", referrer);

  const os = window.navigator.userAgent;
  alert(os);

  const track = async () => {
    const baseUrl = "http://localhost:3000";
    const url = `${baseUrl}/api/v1/analytics_websites/${currentScript.dataset.websiteId}/track`;
    await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        page: location.pathname,
        referrer: referrer || "",
      }),
    });
  };

  const init = () => {
    if (document.readyState === "complete" && !initialized) {
      track();
      initialized = true;
    }
  };

  document.addEventListener("readystatechange", init, true);

  init();
})(window);
