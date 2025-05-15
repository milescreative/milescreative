//just to test the store in a component


// import csrfStore from '../../data/store'
// import * as cache from '../../data/cache'
import * as myAtom from '../../data/atom'
export default function SomeTest() {
  // const csrfToken = csrfStore.get()
  // const csrfToken2 = cache.get()
  const csrfToken3 = myAtom.csrfStore.get()
  return<div>
    {/* <div>{`csrfToken: ${csrfToken}`}</div> */}
    {/* <div>{`csrfToken2: ${csrfToken2}`}</div> */}
    <div>{`csrfToken3: ${csrfToken3}`}</div>
  </div>
}
