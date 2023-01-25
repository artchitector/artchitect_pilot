<template>
  <section>
    <p class="has-text-centered">Artchitect's heart</p>
    <p>Version: {{ artist.version }}</p>
    <p>Seed: {{ artist.seed }}</p>
    <p>
      Tags count:
      <span v-if="artist.totalTags === 0">-</span>
      <span v-else>{{artist.tags.length}}/{{artist.totalTags}}</span>
    </p>
    <p>Tags: {{ artist.tags.join(', ') }}</p>
    <p>
      Progress:
      <span v-if="artist.currentCardPaintTime === null">
        -
      </span>
      <span v-else>
        {{artist.currentCardPaintTime}}/{{artist.lastCardPaintTime}}
      </span>
    </p>
  </section>
</template>
<script>
export default {
  data() {
    return {
      connection: null,
      artist: {
        version: null,
        seed: null,
        totalTags: 0,
        tags: [],
        currentCardPaintTime: null,
        lastCardPaintTime: null
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
    },
    reset() {
      this.artist.version = null
      this.artist.seed = null
      this.artist.tags = []
      this.artist.totalTags = 0
      this.artist.currentCardPaintTime = null
      this.artist.lastCardPaintTime = null
    }
  },
  mounted() {
    const self = this
    if (process.server === true) {
      return
    }
    console.log("❤️: Starting connection to WebSocket Server on ", process.env.WS_URL)
    this.connection = new WebSocket(process.env.WS_URL)

    this.connection.onerror = function (error) {
      console.log(error)
    }

    this.connection.onmessage = function (event) {
      event = JSON.parse(event.data);
      if (event.Name === 'artist') { // card is in work now
        let artistState = JSON.parse(event.Payload)
        console.log("❤️:", artistState)
        if (!artistState.Version) {
          self.reset()
        } else {
          self.updateState(artistState)
        }
      }
    }

    this.connection.onopen = function (event) {
      console.log("Successfully connected to the echo websocket server...")
    }
  },
  beforeDestroy() {
    this.connection.close()
  },
}
</script>
