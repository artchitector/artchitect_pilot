<template>
  <div>
      <textarea class="textarea" v-model="pray" :disabled="locked"
                :placeholder="$t('textarea_placeholder')"></textarea>
    <div class="has-text-centered mt-2">
      <input type="checkbox" v-model="wish" :disabled="locked"> {{$t('wish')}}<br/>
      <button class="button" @click="clear()" :disabled="locked">{{$t('amen')}}</button>
    </div>
    <div class="notification is-primary is-light has-text-centered mt-2" v-if="loading || card_id">
      Пока ваша картина создаётся, вы можете <a href="#" @click.prevent="isDonateVisible = true">пожертвовать на храм</a>
    </div>
    <div class="notification is-info is-light has-text-centered" v-if="loading">
      {{$t('good_time_for_pray')}}<br/>
      {{$t('attempt')}}: {{ attempt }}/{{ maxAttempts }}{{ loader }}<br/>
      <span class="is-size-7">{{$t('usually_time')}}</span>
    </div>
    <div class="notification mt-4 is-danger" v-if="error">
      {{$t('something_wrong')}} - {{ error }}
    </div>
    <div class="image-container has-text-centered" v-if="card_id">
      {{$t('answer')}}
      <br/>
      <a :href="`/card/${card_id}`" target="_blank">
        <img :src="`/api/image/m/${card_id}`"/>
      </a>
    </div>
    <div class="has-text-centered mt-3">
      <donate :isVisible="isDonateVisible" @close="isDonateVisible = false"/>
      <a href="#" @click.prevent="isDonateVisible = true">пожертвовать на храм</a>
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
      loader: ".",
      isDonateVisible: false
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
