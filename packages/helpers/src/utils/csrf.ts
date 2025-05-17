/* eslint-disable prefer-const */

import { store } from "./store";


export const createCsrfStore = (csrfUrl: string, csrfHeader: string) =>{
  const csrfStore = store<string | null>(null);
  csrfStore.onMount(() => {
    fetch(csrfUrl, {
      credentials: "include",
      method: "GET",
    })
      .then((res) => res.headers.get(csrfHeader))
      .then((data) => csrfStore.set(data));
  });
  csrfStore.eager()
  return csrfStore;
}

const defaultCsrfUrl = 'http://localhost:3000/api/auth/csrf'
const defaultCsrfHeader = 'X-CSRF-Token'
const defaultCsrfStore = createCsrfStore(defaultCsrfUrl, defaultCsrfHeader)


type RequestInfo = Request | string | URL;
type RequestInitOmit = (Omit<RequestInit, 'credentials'> & {credentials?:'same-origin' | 'include'}) | undefined;

export const safeFetch = (info: RequestInfo, init: RequestInitOmit, getCsrfToken?: () => string | null) => {
  const csrfGetter = getCsrfToken ?? (() => defaultCsrfStore.get() ?? '');
  const csrfToken = csrfGetter();


  let {headers, credentials, method, ...rest} = init ?? {};
  if (csrfToken) {
    headers = {...headers, [defaultCsrfHeader]: csrfToken};
  }
  if (headers && 'Content-Type' in headers === false) {
    headers = {...headers, 'Content-Type': 'application/json'};
  }

  credentials = credentials ?? 'include';




  const requestInit = {...rest, headers, credentials, method: method ?? 'POST'};
  return fetch(info, requestInit)
}

