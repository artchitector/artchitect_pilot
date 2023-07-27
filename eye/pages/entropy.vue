<i18n>
{
  "en": {
    "title": "Artchitect - entropy"
  },
  "ru": {
    "title": "Artchitect - энтропия"
  }
}
</i18n>

<template>
  <section>
    <div class="has-text-centered">
      <h1>Source</h1>
      <img v-if="images.source !== null" :src="`data:image/jpeg;base64, ${images.source}`"/>
    </div>
  </section>
</template>

<script>

import WsConnection from "~/utils/ws_connection";

export default {
  name: "entropy",
  data() {
    return {
      player: null,
      logPrefix: "🎆",
      status: {
        error: null,
        reconnecting: null,
      },
      maintenance: false,
      connection: null,
      images: {
        source: null
      }
    }
  },
  head() {
    return {
      title: this.$t('title')
    }
  },
  mounted() {
    if (process.env.SOUL_MAINTENANCE === 'true') {
      this.maintenance = true
      return
    }
    this.connection = new WsConnection(process.env.WS_URL, this.logPrefix, ['entropy'], 100)
    this.connection.onmessage((channel, message) => {
      this.status.error = null
      this.status.reconnecting = null
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
      console.log(`${this.logPrefix}: connection established`)
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
    onMessage(chan, msg) {
      if (!msg.Image) {
        return
      }

      switch (msg.Phase) {
        case 'source':
          console.log(`${this.logPrefix}: update image ${msg.Phase}`)
          this.images.source = msg.Image
          break
        default:
          alert(`new phase ${msg.Phase}`)
      }
    }
  }
}
</script>


<style scoped lang="scss">

</style>
