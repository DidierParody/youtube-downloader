import { c as create_ssr_component, s as subscribe, d as add_attribute, f as each, e as escape } from "../../chunks/ssr.js";
import "@sveltejs/kit/internal";
import "../../chunks/exports.js";
import "../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../chunks/state.svelte.js";
import { t as token } from "../../chunks/stores.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $$unsubscribe_token;
  $$unsubscribe_token = subscribe(token, (value) => value);
  let downloads = [];
  let url = "";
  $$unsubscribe_token();
  return `${$$result.head += `<!-- HEAD_svelte-1gfiufh_START -->${$$result.title = `<title>Dashboard - YouTube Downloader</title>`, ""}<!-- HEAD_svelte-1gfiufh_END -->`, ""} <h2 data-svelte-h="svelte-1fh764o">Your Downloads</h2> <form class="form"><h3 data-svelte-h="svelte-16jtag3">Add New Download</h3> <input type="text" placeholder="YouTube URL" required${add_attribute("value", url)}> <button type="submit" class="btn" data-svelte-h="svelte-191embs">Download</button></form> ${``} <div>${downloads.length ? each(downloads, (dl) => {
    return `<div class="card"><h4>${escape(dl.title || "Untitled")}</h4> <p>URL: ${escape(dl.url)}</p> <p>Status: <strong>${escape(dl.status)}</strong></p> <p>Created: ${escape(new Date(dl.created_at).toLocaleString())}</p> </div>`;
  }) : `<p data-svelte-h="svelte-xhew2m">No downloads yet.</p>`}</div>`;
});
export {
  Page as default
};
