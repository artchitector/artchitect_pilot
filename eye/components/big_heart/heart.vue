<template>
  <div>
    <creation ref="main"/>
  </div>
</template>

<script>
import WsConnection from "@/utils/ws_connection";
import Creation from "@/components/big_heart/layout/creation.vue";
export default {
  name: "heart",
  components: {Creation},
  data() {
    return {
      logPrefix: '❤️',
      status: {
        error: null,
        reconnecting: null,
      },
      maintenance: false,
      connection: null,
      // message: null, // message is simple array ["channelName", {message body...}]
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(process.env.WS_URL, this.logPrefix, ['creation', 'lottery', 'unity', 'heart'], 100)
    this.connection.onmessage((channel, message) => {
      this.status.error = null
      this.status.reconnecting = null
      // this.message = [channel, message]
      this.onMessage(channel, message)
    })
    this.connection.onerror((err) => {
      this.status.error = err
    })
    this.connection.onreconnecting((attempt, maxAttempts) => {
      console.log(`${this.logPrefix}: RECONNECTING ${attempt}/${maxAttempts}`)
      this.status.reconnecting = {attempt, maxAttempts}
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
    })
    this.connection.onopen(() => {
      this.status.reconnecting = null
      this.status.error = null
      console.log(`${this.logPrefix}: connection established 🍏`)
    })
    this.connection.connect()
  },
  beforeDestroy() {
    if (!this.maintenance) {
      this.connection.close()
      this.connection = null
    }
  },
  methods: {
    onMessage (channelName, message) {
      // В Сердце может находиться не тот компонент, по которому пришло новое сообщение.
      // Такое бывает, когда режим Архитектора переключается на иную задачу
      // (например, нарисовал и пошёл собирать множество)
      console.log(`${this.logPrefix}: new message`, `channel:${channelName}`, message)
      this.$refs.main.onMessage(channelName, message)
    },
  }
}
</script>

<style scoped>

</style>
