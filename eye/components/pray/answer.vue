<template>
  <div class="has-text-centered">
    <template v-if="loading">
      <loader v-if="loading"/>
      <p>Your answer is loading</p>
      <template v-if="state">
        <template v-if="state.State === 'waiting'">
          <p v-if="state.Queue === 0">Your request is next in queue. Time to wait - less than 60 seconds.</p>
          <p v-else-if="state.Queue > 0">Requests in queue before your: {{ state.Queue }}. Time to wait -
            {{ (state.Queue + 1) * 30 + 30 }} seconds.</p>
        </template>
        <template v-else-if="state.State === 'running'">
          <p>Your answer is in work. Time to wait - less than 30 seconds</p>
        </template>
      </template>
    </template>
    <div v-else-if="error" class="notification is-danger">{{ error }}</div>
    <div v-else-if="state && state.State === 'answered' && state.Answer > 0">
      <p>Your answer is:</p>
      <NuxtLink :to="`/card/${state.Answer}`">
        <img :src="`/api/image/f/${state.Answer}`"/>
      </NuxtLink>
    </div>
  </div>
</template>

<script>

export default {

  name: "answer",
  props: ['id'],
  data() {
    return {
      loading: true,
      error: null,
      state: null,
      interval: null,
    }
  },
  async mounted() {
    const id = parseInt(this.id)
    if (!id) {
      this.error = "ID must be type uint"
      return
    }
    const password = localStorage.getItem("last_pray_password")
    if (!password) {
      this.error = "Access to your request was lost. Please, make new request. Sorry :("
      return
    }
    this.startListening(id, password)
  },
  methods: {
    async startListening(id, password) {
      this.loading = true;
      this.interval = setInterval(async () => {
        try {
          this.state = await this.$axios.$post(`/pray/answer`, {
            id: id,
            password: password,
          })
          console.log(`🙏: current pray state ${this.state.State}`)
          if (this.state.State === "answered") {
            this.loading = false;
            clearInterval(this.interval)
          }
        } catch (e) {
          this.error = e.message;
          this.loading = false;
          clearInterval(this.interval)
        }
      }, 1000)
    }
  }
}
</script>

<style scoped>

</style>
