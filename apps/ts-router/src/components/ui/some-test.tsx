//just to test the store in a component


import csrfStore from '../../data/store'
import * as cache from '../../data/cache'

export default function SomeTest() {
  const csrfToken = csrfStore.get()
  const csrfToken2 = cache.get()
  return<div>
    <div>{`csrfToken: ${csrfToken}`}</div>
    <div>{`csrfToken2: ${csrfToken2}`}</div>
  </div>
}
