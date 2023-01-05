<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h1 class="mt-5">Login</h1>
        <hr />
        <form-tag @myevent="submitHandler" name="myform" event="myevent">
          <text-input v-model="email" label="Email" type="email" name="email" required="true"> </text-input>

          <text-input v-model="password" label="Password" type="password" name="password" required="true">
          </text-input>

          <hr />
          <input type="submit" class="btn btn-primary" value="Login" />
        </form-tag>
      </div>
    </div>
  </div>
</template>

<script>
import FormTag from './forms/FormTag.vue'
import TextInput from './forms/TextInput.vue'
import { store } from '../store/'
import router from '../router/'
import notie from 'notie'

export default {
  name: 'login',
  components: {
    FormTag,
    TextInput,
  },
  data() {
    return {
      email: '',
      password: '',
      store,
    }
  },
  methods: {
    submitHandler() {
      console.log('submitHandler called - success!')

      const payload = {
        email: this.email,
        password: this.password,
      }

      const requestOptions = {
        method: 'POST',
        body: JSON.stringify(payload),
      }

      fetch('http://localhost:8081/users/login', requestOptions)
        .then((response) => response.json())
        .then((res) => {
          if (res.error) {
            console.log('Error:', res.message)
            notie.alert({ type: 'error', text: res.message })
          } else {
            this.store.token = res.data.token.token
            console.log(res.data.token.token)
            router.push('/')
          }
        })
    },
  },
}
</script>
