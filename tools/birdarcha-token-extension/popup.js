async function render() {
  const { lastStatus, lastStatusAt } = await chrome.storage.local.get([
    "lastStatus",
    "lastStatusAt",
  ]);
  document.getElementById("status").textContent = lastStatus
    ? `${lastStatusAt}\n${lastStatus}`
    : "no runs yet";
}

document.getElementById("refresh").addEventListener("click", async () => {
  document.getElementById("status").textContent = "refreshing…";
  await chrome.runtime.sendMessage({ type: "refresh-now" });
  setTimeout(render, 1500);
});

document.getElementById("options").addEventListener("click", () => {
  chrome.runtime.openOptionsPage();
});

render();
