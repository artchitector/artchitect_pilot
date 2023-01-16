<template>
  <div>
    <viewer ref="viewer"/>
    <h3 class="is-size-4 has-text-centered mb-4" v-if="count > 0">{{$t('last')}} {{count}} {{$t('cards')}}</h3>
    <h3 class="is-size-4"v-else>{{$t('last')}} {{$t('cards')}}</h3>
    <div class="notification is-primary" v-if="!cards.length && $fetchState.pending">
      {{$t('loading')}}...
    </div>
    <div class="notification is-danger" v-else-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>

    <div v-else>
      <div class="columns" v-for="line in lines">
        <div class="column" v-for="artwork in line">
          <artwork-view :artwork="artwork" @select="select(artwork.ID)"/>
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
  },
  methods: {
    select(id) {
      const ids = []
      this.cards.forEach((card) => {
        ids.push(card.ID)
      })
      this.$refs.viewer.show(ids, id);
    }
  }
}
</script>
