const $ = (id) => document.getElementById(id);

async function load() {
  const { apiBase, apiKey, intervalMin, lastStatus, lastStatusAt } =
    await chrome.storage.local.get([
      "apiBase",
      "apiKey",
      "intervalMin",
      "lastStatus",
      "lastStatusAt",
    ]);
  $("apiBase").value = apiBase || "";
  $("apiKey").value = apiKey || "";
  $("intervalMin").value = intervalMin || 180;
  if (lastStatus) {
    $("status").textContent = `${lastStatusAt || ""} — ${lastStatus}`;
  }
}

$("save").addEventListener("click", async () => {
  await chrome.storage.local.set({
    apiBase: $("apiBase").value.trim(),
    apiKey: $("apiKey").value.trim(),
    intervalMin: Math.max(5, Number($("intervalMin").value) || 180),
  });
  await chrome.runtime.sendMessage({ type: "reschedule" });
  $("status").textContent = "saved & rescheduled";
});

$("refresh").addEventListener("click", async () => {
  $("status").textContent = "refreshing…";
  await chrome.runtime.sendMessage({ type: "refresh-now" });
  setTimeout(load, 1500);
});

load();
