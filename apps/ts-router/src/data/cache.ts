// csrfStore.js
let value:string | null = null;
let listeners:any[] = [];
let mounted = false;

export function get() {
  return value;
}

export function set(newVal:string | null) {
  value = newVal;
  listeners.forEach((fn) => fn(value));
}

export function listen(fn:any) {
  listeners.push(fn);
  return () => {
    listeners = listeners.filter((f) => f !== fn);
  };
}

function mount() {
  if (mounted) return;
  mounted = true;
  fetch('http://localhost:3000/api/auth/csrf', {
    credentials: 'include',
    method: 'GET',
  })
    .then((res) => res.headers.get('X-Csrf-Token'))
    .then(set);
}

mount(); // auto-fetch once on import

// listen((val:string | null) => {
//   console.log('New CSRF token test:', val);
// });
