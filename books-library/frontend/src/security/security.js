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
    headers.append('Authorization', 'Bearer ' + store.token)

    return {
      method: 'POST',
      body: JSON.stringify(payload),
      headers: headers,
    }
  },

  checkToken: function () {
    if (store.token !== '') {
      const payload = {
        token: store.token,
      }

      const headers = new Headers()
      headers.append('Content-Type', 'application/json')

      let requestOptions = {
        method: 'POST',
        body: JSON.stringify(payload),
        headers: headers,
      }

      fetch(process.env.VUE_APP_API_URL + '/validate-token', requestOptions)
        .then((response) => response.json())
        .then((data) => {
          if (data.error) {
            console.log(data.error)
          } else {
            if (!data.data) {
              store.token = ''
              ;(store.user = {}),
                (document.cookie =
                  '_site_data=; Path=/; ' +
                  'SameSite=strict; Secure; ' +
                  'Expires=Thu, 01 Jan 1970 00:00:01 GMT;')
            }
          }
        })
    }
  },
}

export default Security
