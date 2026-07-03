import { c as create_ssr_component, s as subscribe, e as escape } from "../../chunks/ssr.js";
import "@sveltejs/kit/internal";
import "../../chunks/exports.js";
import "../../chunks/utils2.js";
import "@sveltejs/kit/internal/server";
import "../../chunks/state.svelte.js";
import { u as user } from "../../chunks/stores.js";
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $user, $$unsubscribe_user;
  $$unsubscribe_user = subscribe(user, (value) => $user = value);
  $$unsubscribe_user();
  return `<header><h1 data-svelte-h="svelte-1qg24lm">🎬 YouTube Downloader</h1> ${$user ? `<button class="btn">Logout (${escape($user.username)})</button>` : ``}</header> <main class="container">${slots.default ? slots.default({}) : ``}</main>`;
});
export {
  Layout as default
};
