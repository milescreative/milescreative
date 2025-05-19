import { createCsrfStore } from '../../../../packages/backend-sdk/src'

export const csrfStore = createCsrfStore('http://localhost:3000/api/auth/csrf', 'X-Csrf-Token')
