<template>
  <div class="container">
    <div class="row">
      <div class="col-md-2">
        <img
          v-if="this.ready"
          class="img-fluid img-thumbnail"
          :src="`${imgPath}/covers/${book.slug}.jpg`"
          alt="cover" />
      </div>
      <div class="col-md-10">
        <template v-if="this.ready">
          <h3 class="mt-3">{{ book.title }}</h3>
          <hr />
          <p>
            <strong>Author:</strong> {{ book.author.author_name }}<br />
            <strong>Published:</strong> {{ book.publication_year }}
          </p>
          <p>
            {{ book.description }}
          </p>
        </template>
        <p v-else>Loading...</p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      book: {},
      imgPath: process.env.VUE_APP_IMAGE_URL,
      ready: false,
    }
  },
  activated() {
    // activated is fired anytime a cached component is redisplayed, and anytime mounted called
    fetch(process.env.VUE_APP_API_URL + '/books/' + this.$route.params.bookName)
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          this.$emit('error', data.message)
        } else {
          this.book = data.data
          this.ready = true
          console.log('Title is', this.book.title)
        }
      })
  },
  deactivated() {
    this.ready = false
  },
}
</script>
