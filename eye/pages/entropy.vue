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
      <img v-if="images.source !== null" :src="`data:image/jpeg;base64, ${images.source}`"
           alt="loading source stream"/>
    </div>
    <div class="has-text-centered">
      <h1>Noise</h1>
      <img v-if="images.noise !== null" :src="`data:image/jpeg;base64, ${images.noise}`" style="height: 256px; width: 256px;"
           alt="loading noise stream"/>
    </div>
    <div class="has-text-centered">
      <h1>Shrink</h1>
      <img v-if="images.shrink !== null" :src="`data:image/png;base64, ${images.shrink}`"
           style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading shrink stream"/>
    </div>
    <div class="has-text-centered">
      <h1>Bytes</h1>
      <img v-if="images.bytes !== null" :src="`data:image/png;base64, ${images.bytes}`"
           style="width: 64px; height: 64px; image-rendering: pixelated;" alt="loading bytes stream"/>
      <br/>
      <span class="is-size-7" v-html="entropy.bytes ? entropy.bytes.match(/.{1,8}/g).join('<br/>') : '-'"></span><br/>
      <h1>Resulting value (scale 0 to 1)</h1>
      <span>{{ entropy.float }}</span><br/>
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
        source: null,
        noise: null,
        shrink: null,
        bytes: null
      },
      entropy: {
        bytes: "",
        float: 0.0
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

      if (msg.EntropyAnswer > 0) {
        this.entropy.bytes = msg.EntropyAnswerByte
        this.entropy.float = msg.EntropyAnswer
      }

      switch (msg.Phase) {
        case 'source':
          this.images.source = msg.Image
          break
        case 'noise':
          this.images.noise = msg.Image
          break
        case 'shrink':
          this.images.shrink = msg.Image
          break
        case 'bytes':
          this.images.bytes = msg.Image
          break
        default:
          console.warn(`new phase ${msg.Phase}`)
      }
    }
  }
}
</script>


<style scoped lang="scss">

</style>
