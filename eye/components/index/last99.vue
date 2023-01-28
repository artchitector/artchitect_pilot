<template>
  <div>
    <h3 class="is-size-4 has-text-centered mb-4">last 99 cards</h3>
    <div v-if="$fetchState.pending" class="notification has-text-centered">
      loading...
    </div>
    <div v-else-if="$fetchState.error" class="notification is-danger">
      {{ $fetchState.error.message }}
    </div>
    <div v-else-if="!this.cards || !this.cards.length" class="notification is-danger">
      cards not loaded :(
    </div>
    <div v-else>
      <div v-if="wsStatus.error" class="notification is-warning is-size-7 has-text-centered">
        websocket listening error: {{wsStatus.error.message}}
      </div>
      <div v-else-if="wsStatus.reconnecting" class="notification is-size-7 has-text-centered">
        websocket connecting {{wsStatus.reconnecting.attempt}}/{{wsStatus.reconnecting.maxAttempts}}
      </div>
      <cardlist :cards="this.cards" cards-in-column="3" card-size="m" visible-count="33"/>
    </div>
  </div>
</template>

<script>
import Cardlist from "@/components/list/cardlist.vue";
import WsConnection from "@/utils/ws_connection";

export default {
  name: "last99",
  components: {Cardlist},
  data() {
    return {
      connection: null,
      wsStatus: {
        error: null,
        reconnecting: null,
      },
      cards: []
    }
  },
  mounted() {
    this.connection = new WsConnection(process.env.WS_URL, '🌄', ['new_card'], 10)
    this.connection.onmessage((channel, newCard) => {
      this.wsStatus.error = null;
      this.wsStatus.reconnecting = null;
      const removedCard = this.cards[this.cards.length - 1];
      const cards = this.cards.slice(0, this.cards.length - 1)
      cards.unshift(newCard)
      this.cards = cards;
      console.log(`🌄: new card (id=${newCard.ID}), removed (id=${removedCard.ID})`,)
    })
    this.connection.onerror((err) => this.wsStatus.error = err)
    this.connection.onreconnecting((attempt, maxAttempts) => this.wsStatus.reconnecting = {
      attempt: attempt,
      maxAttempts: maxAttempts
    })
    this.connection.onopen(() => {
      this.wsStatus.error = null
      this.wsStatus.reconnecting = null
    })
    this.connection.connect()
  },
  beforeDestroy() {
    this.connection.close()
    this.connection = null;
  },
  async fetch() {
    try {
      this.cards = await this.$axios.$get('/last_paintings/99')
      console.log('[last99] loaded cards', this.cards.length)
    } catch (e) {
      if (this.connection) {
        this.connection.close()
      }
      throw e;
    }
  }
}
</script>

<style scoped>

</style>
