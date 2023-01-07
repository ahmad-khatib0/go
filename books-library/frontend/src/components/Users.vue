<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <h1 class="mt-3">All Users</h1>
      </div>
    </div>
  </div>
</template>

<script>
import Secority from '../security/security'

export default {
  beforeMount() {
    Secority.requireAuth()
  },
}
</script>
