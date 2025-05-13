import { atom, onMount } from 'nanostores'

const csrfStore = atom<any | null>(null)

onMount(csrfStore, () => {
  fetch('http://localhost:3000/api/auth/csrf', {
    credentials: 'include',
    method: 'GET',
  })
    .then(res => res.headers.get('X-Csrf-Token'))
    .then(data => csrfStore.set(data))
})

csrfStore.listen(() => {
  const data = csrfStore.get()
  console.log('New data:', data)
})

export default csrfStore
