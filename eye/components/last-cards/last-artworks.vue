<template>
  <div>
    <a id="cards"></a>
    <h3 class="is-size-4 has-text-centered mb-4" v-if="count > 0">last {{count}} cards</h3>
    <h3 class="is-size-4"v-else>last cards</h3>
    <div class="notification is-primary" v-if="!cards.length && $fetchState.pending">
      loading...
    </div>
    <div class="notification is-danger" v-else-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>

    <div v-else>
      <div class="has-text-centered mb-4">
          <nuxt-link v-for="page in pages" :to="page.url">
            {{page.caption}}
          </nuxt-link>
      </div>
      <div class="columns" v-for="line in lines">
        <div class="column" v-for="artwork in line">
          <artwork-view :artwork="artwork"/>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import ArtworkView from "~/components/last-cards/artwork-view.vue";

export default {
  components: {ArtworkView},
  data() {
    return {
      cards: [],
      error: null
    }
  },
  async fetch() {
    this.cards = await this.$axios.$get('/last_paintings/100')
    console.log('loaded', this.cards.length)
  },
  computed: {
    pages() {
      if (!this.cards.length) {
        return []
      }
      const lastId = this.cards[0].ID
      const pages = []
      let from, to
      for (let i = 0; i < 5; i++) {
        from = lastId - (i * 100) - 1
        to = from - 100
        pages.push({
          "url": `/list?from=${from}&to=${to}`,
          "caption": `(page${i+1}:${from}-${to})`
        })
      }
      return pages
    },
    count() {
      return this.cards.length
    },
    lines() {
      if (!this.cards.length) {
        return [];
      }
      const lines = [
        []
      ]
      for (let i = 0; i < this.cards.length; i++) {
        const artwork = this.cards[i]
        let lastLine = lines.length - 1
        if (lines[lastLine].length < 3) {
          lines[lastLine].push(artwork)
        } else {
          lines.push([artwork])
        }
      }
      return lines
    }
  }
}
</script>
