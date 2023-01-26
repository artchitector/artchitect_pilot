<template>
  <section class="heart">
    <viewer ref="viewer"/>
    <h3 class="has-text-centered is-size-4">Artchitect's heart</h3>
    <hr class="divider"/>
    <div class="notification is-danger" v-if="status.error">
      connection failed
    </div>
    <div v-else-if="!status.connected" class="has-text-centered">
      connecting {{this.status.attempts}}/{{this.status.maxAttempts}}
    </div>
    <div v-else-if="status.artchitectShutdown" class="notification is-warning">
      Artchitect is offline now
    </div>
    <div v-else class="is-flex is-flex-direction-row">
      <div class="image-container">
        <img v-if="!artist.cardId" src="/in-progress.jpeg"/>
        <a v-else :href="`/card/${artist.cardId}`" target="_blank" @click.prevent="viewer()">
          <img :src="`/api/image/s/${artist.cardId}`"/>
        </a>
      </div>
      <div class="info-container">
        <p>Version: {{ artist.version }}</p>
        <p>Seed: {{ artist.seed }}</p>
        <p>
          Tags count:
          <span v-if="artist.totalTags === 0">-</span>
          <span v-else>{{ artist.tags.length }}/{{ artist.totalTags }}</span>
        </p>
        <p class="is-size-7">Tags: {{ artist.tags.join(', ') }}</p>
        <p>
          Paint progress:
          <span v-if="artist.currentCardPaintTime === null">
        -
      </span>
          <span v-else>
        {{ artist.currentCardPaintTime }}/{{ artist.lastCardPaintTime }}
      </span>
        </p>
        <p>
          Enjoy progress:
          <span v-if="artist.currentEnjoyTime === null">
          -
          </span>
          <span v-else>
            {{ artist.currentEnjoyTime }}/{{ artist.totalEnjoyTime }}
          </span>
        </p>
      </div>
    </div>
  </section>
</template>
<script>
import connection from './connection'
export default {
  data() {
    return {
      connection: null,
      status: {
        connected: false,
        error: null,
        attempts: 1,
        maxAttempts: 100,
        lastMessageTime: null,
        artchitectShutdown: false,
      },
      artist: {
        version: null,
        seed: null,
        totalTags: 0,
        tags: [],
        currentCardPaintTime: null,
        lastCardPaintTime: null,
        totalEnjoyTime: null,
        currentEnjoyTime: null,
        cardId: null,
      }
    }
  },
  methods: {
    updateState(state) {
      if (!!state.Version) {
        this.artist.version = state.Version
      }
      if (!!state.Seed) {
        this.artist.seed = state.Seed
      }
      if (state.TagsCount > 0) {
        this.artist.totalTags = state.TagsCount
      }
      if (!!state.Tags && state.Tags.length) {
        this.artist.tags = state.Tags
      }
      if (!!state.LastCardPaintTime) {
        this.artist.lastCardPaintTime = state.LastCardPaintTime
      }
      if (!!state.CurrentCardPaintTime) {
        this.artist.currentCardPaintTime = state.CurrentCardPaintTime
      }
      if (!!state.CurrentEnjoyTime) {
        this.artist.currentEnjoyTime = state.CurrentEnjoyTime
      }
      if (!!state.EnjoyTime) {
        this.artist.totalEnjoyTime = state.EnjoyTime
      }
      if (!!state.CardID) {
        this.artist.cardId = state.CardID
      }
    },
    reset() {
      this.artist.version = null
      this.artist.seed = null
      this.artist.tags = []
      this.artist.totalTags = 0
      this.artist.currentCardPaintTime = null
      this.artist.lastCardPaintTime = null
      this.artist.currentEnjoyTime = null
      this.artist.totalEnjoyTime = null
      this.artist.cardId = null
    },
    viewer() {
      const ids = [this.artist.cardId]
      this.$refs.viewer.show(ids, this.artist.cardId);
    }
  },
  mounted() {
    this.connection = connection

    setInterval(() => {
      if (this.status.artchitectShutdown) {
        return
      }
      if (this.status.lastMessageTime === null) {
        this.status.artchitectShutdown = true
        return
      }
      if (new Date().getTime() - this.status.lastMessageTime.getTime() > 5000) {
        this.status.artchitectShutdown = true
        return
      }
    }, 5000)

    this.connection.onmessage((state) => {
      console.log('onmessage', state)
      this.status.lastMessageTime = new Date()
      this.status.artchitectShutdown = false
      if (!state.Version) {
        this.reset()
      } else {
        this.updateState(state)
      }
    })
    this.connection.onerror((err) => {
      console.log('onerror', err)
    })
    this.connection.onconnect(() => {
      console.log('onconnect')
      this.status.attempts = 1
      this.status.connected = true
    })
    this.connection.onclose(() => {
      this.status.connected = false
      setTimeout(() => {
        this.status.attempts += 1
        if (this.status.attempts < this.status.maxAttempts) {
          this.connection.connect(process.env.WS_URL)
        } else {
          this.status.error = "connection failed"
        }
      }, 1000)
    })
    this.connection.connect(process.env.WS_URL)

  },
  beforeDestroy() {
    this.connection.close()
  },
}
</script>
<style lang="scss">
.heart div.image-container {
  min-width: 170px;
  width: 170px;
  padding-right: 10px;

  a {
    display: block;
  }
}

.heart hr.divider {
  margin: 0 0 0.5rem 0;
}

.heart div.info-container {

}
</style>
