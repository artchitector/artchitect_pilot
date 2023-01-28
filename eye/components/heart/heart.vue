<template>
  <section class="heart">

    <h3 class="has-text-centered is-size-4">Artchitect's heart</h3>
    <hr class="divider"/>
    <div v-if="status.error" class="notification is-danger has-text-centered">
      Error: {{status.error.message}}
    </div>
    <div v-else-if="status.reconnecting" class="notification has-text-centered">
      connecting {{status.reconnecting.attempt}}/{{status.reconnecting.maxAttempts}}
    </div>
    <p v-else-if="!stateChannel">
      no state
    </p>
    <creation v-else-if="stateChannel === 'creation'" :state="state"/>
    <p v-else>
      unknown state {{ stateChannel }}
    </p>
  </section>
</template>
<script>
import WsConnection from '../../utils/ws_connection'
import Creation from "@/components/heart/layout/creation.vue";

export default {
  components: {Creation},
  data() {
    return {
      status: {
        error: null,
        reconnecting: null,
      },
      connection: null,
      stateChannel: null,
      state: null,
    }
  },
  mounted() {
    this.connection = new WsConnection(process.env.WS_URL, '🧡', ['creation'], 10)

    this.connection.onmessage((channel, state) => {
      console.log('🧡: new message', channel, state)
      this.status.error = null
      this.status.reconnecting = null
      this.stateChannel = channel
      this.state = state
    })
    this.connection.onerror((err) => {
      this.status.error = err
    })
    this.connection.onreconnecting((attempt, maxAttempts) => {
      console.log(`🧡: RECONNECTING ${attempt}/${maxAttempts}`)
      this.status.reconnecting = {
        attempt: attempt,
        maxAttempts: maxAttempts,
      }
    })
    this.connection.connect()
  },
  beforeDestroy() {
    this.connection.close()
    this.connection = null
  },
}
</script>
<style lang="scss">

</style>
