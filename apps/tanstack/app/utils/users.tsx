export type User = {
  id: number
  name: string
  email: string
}

const PROD_URL = import.meta.env.TANSTACK_PROD_URL || 'https://tanstack.mlcr.us'
// eslint-disable-next-line turbo/no-undeclared-env-vars
export const DEPLOY_URL = import.meta.env.PROD
  ? PROD_URL
  : 'http://localhost:3004'
