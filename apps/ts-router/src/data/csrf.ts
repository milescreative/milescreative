
import { atom } from "./atom";


const csrfStore = atom<any | null>(null);

csrfStore.onMount(() => {
  fetch("http://localhost:3000/api/auth/csrf", {
    credentials: "include",
    method: "GET",
  })
    .then((res) => res.headers.get("X-Csrf-Token"))
    .then((data) => csrfStore.set(data));
});
csrfStore.eager()

const csrfTokenHeader = 'X-CSRF-Token';

type RequestInitOmit = (Omit<RequestInit, 'credentials'> & {credentials?:'same-origin' | 'include'}) | undefined;

export const safeFetch = (info: RequestInfo, init: RequestInitOmit, getCsrfToken?: () => string) => {
  const csrfGetter = getCsrfToken ?? (() => csrfStore.get() ?? '');
  const csrfToken = csrfGetter();

  let {headers, credentials, method, ...rest} = init ?? {};
  headers = {...headers, [csrfTokenHeader]: csrfToken};
  if (headers && 'Content-Type' in headers === false) {
    headers = {...headers, 'Content-Type': 'application/json'};
  }

  credentials = credentials ?? 'include';




  const requestInit = {...rest, headers, credentials, method: method ?? 'POST'};
  return fetch(info, requestInit)
}

