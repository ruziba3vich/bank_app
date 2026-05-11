const ALARM_NAME = "birdarcha-refresh";
const BIRDARCHA_URL = "https://office.birdarcha.uz/";
const DEFAULT_INTERVAL_MIN = 180;
const TOKEN_KEY = "access-token";

async function getConfig() {
  const { apiBase, apiKey, intervalMin } = await chrome.storage.local.get([
    "apiBase",
    "apiKey",
    "intervalMin",
  ]);
  return {
    apiBase: (apiBase || "").replace(/\/$/, ""),
    apiKey: apiKey || "",
    intervalMin: Number(intervalMin) || DEFAULT_INTERVAL_MIN,
  };
}

async function setStatus(text) {
  const ts = new Date().toISOString();
  await chrome.storage.local.set({ lastStatus: text, lastStatusAt: ts });
  console.log(`[birdarcha-refresher] ${ts} ${text}`);
}

async function findBirdarchaTab() {
  const tabs = await chrome.tabs.query({ url: "https://office.birdarcha.uz/*" });
  return tabs[0] || null;
}

async function ensureBirdarchaTab() {
  let tab = await findBirdarchaTab();
  if (tab) return tab;
  tab = await chrome.tabs.create({ url: BIRDARCHA_URL, active: false });
  await new Promise((resolve) => {
    function listener(tabId, info) {
      if (tabId === tab.id && info.status === "complete") {
        chrome.tabs.onUpdated.removeListener(listener);
        resolve();
      }
    }
    chrome.tabs.onUpdated.addListener(listener);
  });
  return tab;
}

async function reloadAndWait(tabId) {
  await chrome.tabs.reload(tabId);
  await new Promise((resolve) => {
    function listener(updatedId, info) {
      if (updatedId === tabId && info.status === "complete") {
        chrome.tabs.onUpdated.removeListener(listener);
        resolve();
      }
    }
    chrome.tabs.onUpdated.addListener(listener);
  });
  await new Promise((r) => setTimeout(r, 1500));
}

async function readTokenFromTab(tabId) {
  const [result] = await chrome.scripting.executeScript({
    target: { tabId },
    func: (key) => localStorage.getItem(key),
    args: [TOKEN_KEY],
  });
  return result?.result || null;
}

function isJwtExpired(jwt) {
  try {
    const payload = JSON.parse(atob(jwt.split(".")[1].replace(/-/g, "+").replace(/_/g, "/")));
    if (!payload.exp) return true;
    return payload.exp * 1000 <= Date.now();
  } catch {
    return true;
  }
}

async function postToken(token) {
  const { apiBase, apiKey } = await getConfig();
  if (!apiBase || !apiKey) {
    throw new Error("apiBase and apiKey must be configured in options");
  }
  const url = `${apiBase}/api/v1/entrepreneurs/birdarcha-token`;
  const resp = await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      "X-Api-Key": apiKey,
    },
    body: JSON.stringify({ token }),
  });
  if (!resp.ok) {
    const body = await resp.text();
    throw new Error(`HTTP ${resp.status}: ${body}`);
  }
  return resp.json().catch(() => ({}));
}

async function refreshOnce() {
  try {
    const tab = await ensureBirdarchaTab();
    let token = await readTokenFromTab(tab.id);

    if (!token || isJwtExpired(token)) {
      await reloadAndWait(tab.id);
      token = await readTokenFromTab(tab.id);
    }

    if (!token) {
      await setStatus("no token in localStorage — are you logged in?");
      return;
    }
    if (isJwtExpired(token)) {
      await setStatus("token still expired after reload — session may need re-login");
      return;
    }

    await postToken(token);
    await setStatus("token pushed to bank-app successfully");
  } catch (err) {
    await setStatus(`error: ${err.message}`);
  }
}

async function scheduleAlarm() {
  const { intervalMin } = await getConfig();
  await chrome.alarms.clear(ALARM_NAME);
  await chrome.alarms.create(ALARM_NAME, {
    delayInMinutes: 1,
    periodInMinutes: intervalMin,
  });
}

chrome.runtime.onInstalled.addListener(scheduleAlarm);
chrome.runtime.onStartup.addListener(scheduleAlarm);

chrome.alarms.onAlarm.addListener((alarm) => {
  if (alarm.name === ALARM_NAME) refreshOnce();
});

chrome.runtime.onMessage.addListener((msg, _sender, sendResponse) => {
  if (msg?.type === "refresh-now") {
    refreshOnce().then(() => sendResponse({ ok: true }));
    return true;
  }
  if (msg?.type === "reschedule") {
    scheduleAlarm().then(() => sendResponse({ ok: true }));
    return true;
  }
});
