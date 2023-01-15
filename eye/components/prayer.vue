<template>
  <div class="box">
    <div v-if="!visible" class="has-text-centered">
      <button class="button is-success" @click="visible=true">To pray!</button>
    </div>
    <div v-else>
      <textarea class="textarea" v-model="pray"
                placeholder="Place for pray message. Write carefully. Between you and God. Secure, data burns and not being send anywhere."></textarea>
      <div class="has-text-centered mt-2">
        <button class="button" @click="clear()">Pray! (burn text)</button>
        <input type="checkbox" v-model="wish"> Wish one non-random card as God's reply
      </div>
      <div class="loader" v-if="loading">
        you card loading
      </div>
      <div class="is-danger" v-if="error">
        Что-то случилось - {{ error }}
      </div>
      <div class="image-container has-text-centered" v-if="card_id">
        <span>Ответ</span>
        <br/>
        <a :href="`/card/${card_id}`" target="_blank">
          <img :src="`/api/image/m/${card_id}`"/>
        </a>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      visible: false,
      pray: "",
      wish: false,
      loading: false,
      error: null,
      card_id: null,
    }
  },
  methods: {
    async clear() {
      this.pray = ""
      this.card_id = null
      if (this.wish) {
        this.wish = false
        try {
          await this.load()
        } catch (e) {
          this.error = e.message
        } finally {
          this.loading = false
        }
      }
    },
    async load() {
      this.card_id = await this.$axios.$get("/api/answer")
    }
  }
}
</script>
