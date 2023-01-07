import { store } from '../store'
import router from '../router'

let Security = {
  requireAuth: function () {
    if (store.token === '') {
      router.push('/login')
      return false
    }
  },
  requestOptions: function (payload) {
    const headers = new Headers()
    headers.append('Content-Type', 'application/json')
    headers.append('Authorization', 'Beare ' + store.token)

    return {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: headers,
    }
  },
}

export default Security
