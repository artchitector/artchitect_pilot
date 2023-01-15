<template>
  <div>
      <textarea class="textarea" v-model="pray"
                placeholder="Type your pray message. Write carefully. Between you and God. Secure, data burns and not being send anywhere."></textarea>
    <div class="has-text-centered mt-2">
      <input type="checkbox" v-model="wish" :disabled="locked"> Wish one non-random card as God's reply<br/>
      <button class="button" @click="clear()" :disabled="locked">Pray! (burn text)</button>
    </div>
    <div class="notification is-info is-light has-text-centered" v-if="loading">
      Your answer loading. Good time for pray!<br/>
      Attempt: {{ attempt }}/{{ maxAttempts }}{{ loader }}<br/>
      <span class="is-size-7">Usually that takes 2 minutes, but if artchitect is very loaded now, you need try once later.</span>
    </div>
    <div class="notification mt-4 is-danger" v-if="error">
      Что-то случилось - {{ error }}
    </div>
    <div class="image-container has-text-centered" v-if="card_id">
      Ответ
      <br/>
      <a :href="`/card/${card_id}`" target="_blank">
        <img :src="`/api/image/m/${card_id}`"/>
      </a>
    </div>
  </div>

</template>
<script>
export default {
  data() {
    return {
      maxAttempts: 60,
      locked: false,
      attempt: 0,
      pray: "",
      wish: false,
      loading: false,
      error: null,
      pray_id: null,
      card_id: null,
      int: null,
      loader: "."
    }
  },
  methods: {
    async clear() {
      this.pray = ""
      this.card_id = null
      this.attempt = 0
      this.pray_id = 0
      this.locked = false
      clearInterval(this.int)
      this.int = null
      this.fastInt = null

      if (this.wish) {
        this.wish = false
        this.answer()
      }
    },
    async answer() {
      try {
        this.loading = true
        this.error = null
        this.locked = true
        this.pray_id = await this.$axios.$get("/answer")
        this.int = setInterval(async () => {
          this.attempt += 1
          if (this.attempt > this.maxAttempts) {
            clearInterval(this.int)
            clearInterval(this.fastInt)
            this.error = "не смог дождаться ответ... :( попробуйте позже..."
            this.locked = false
            this.loading = false
          }
          let answer = await this.$axios.$get(`/answer/${this.pray_id}`)
          if (answer === 0 || answer === "0") {
            // это ок, если будет ждать еще 5 секунд
          } else {
            this.card_id = answer
            clearInterval(this.int)
            clearInterval(this.fastInt)
            this.locked = false
            this.loading = false
          }
        }, 5000)
        this.fastInt = setInterval(() => {
          if (this.loader === ".") {
            this.loader = ".."
          } else if (this.loader === "..") {
            this.loader = "..."
          } else if (this.loader === "...") {
            this.loader = "."
          }
        }, 500)
      } catch (e) {
        this.error = e.message
        clearInterval(this.int)
        clearInterval(this.fastInt)
        this.locked = false
        this.loading = false
      } finally {
      }


    }
  }
}
</script>
