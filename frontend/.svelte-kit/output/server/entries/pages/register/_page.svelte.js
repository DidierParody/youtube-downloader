import { c as create_ssr_component, d as add_attribute } from "../../../chunks/ssr.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/state.svelte.js";
import "../../../chunks/stores.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let username = "";
  let email = "";
  let password = "";
  return `${$$result.head += `<!-- HEAD_svelte-41pyvy_START -->${$$result.title = `<title>Register - YouTube Downloader</title>`, ""}<!-- HEAD_svelte-41pyvy_END -->`, ""} <h2 data-svelte-h="svelte-d4pubn">Register</h2> <form class="form">${``} ${``} <input type="text" placeholder="Username" required${add_attribute("value", username)}> <input type="email" placeholder="Email" required${add_attribute("value", email)}> <input type="password" placeholder="Password" required${add_attribute("value", password)}> <button type="submit" class="btn" data-svelte-h="svelte-1trnnlp">Register</button> <p data-svelte-h="svelte-mcwoq1">Already have an account? <a href="/login">Login</a></p></form>`;
});
export {
  Page as default
};
