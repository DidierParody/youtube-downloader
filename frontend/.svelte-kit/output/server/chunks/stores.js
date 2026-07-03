import { w as writable } from "./index.js";
const user = writable(null);
const token = writable(typeof window !== "undefined" ? localStorage.getItem("token") : null);
user.subscribe((u) => {
  if (typeof window !== "undefined" && u) {
    localStorage.setItem("user", JSON.stringify(u));
  }
});
token.subscribe((t) => {
  if (typeof window !== "undefined" && t) {
    localStorage.setItem("token", t);
  }
});
export {
  token as t,
  user as u
};
