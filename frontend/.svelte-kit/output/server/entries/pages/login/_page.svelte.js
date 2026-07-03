import { c as create_ssr_component, d as add_attribute } from "../../../chunks/ssr.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/state.svelte.js";
import "../../../chunks/stores.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let email = "";
  let password = "";
  return `${$$result.head += `<!-- HEAD_svelte-mxdx92_START -->${$$result.title = `<title>Login - YouTube Downloader</title>`, ""}<!-- HEAD_svelte-mxdx92_END -->`, ""} <h2 data-svelte-h="svelte-bhb3ah">Login</h2> <form class="form">${``} <input type="email" placeholder="Email" required${add_attribute("value", email)}> <input type="password" placeholder="Password" required${add_attribute("value", password)}> <button type="submit" class="btn" data-svelte-h="svelte-1nmopzr">Login</button> <p data-svelte-h="svelte-nbf4t">Don&#39;t have an account? <a href="/register">Register</a></p></form>`;
});
export {
  Page as default
};
