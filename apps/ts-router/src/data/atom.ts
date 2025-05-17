import { createCsrfStore } from '@milescreative/helpers'

export const csrfStore = createCsrfStore('http://localhost:3000/api/auth/csrf', 'X-Csrf-Token')
