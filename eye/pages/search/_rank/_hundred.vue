<i18n>
{
  "en": {
    "title": "Artchitect - search",
    "subtitle": "search"
  },
  "ru": {
    "title": "Artchitect - поиск",
    "subtitle": "поиск"
  }
}
</i18n>

<template>
  <section>
    <h3 class="has-text-centered is-size-4">{{ subtitle }}</h3>
    <template v-if="$fetchState.pending">
      <div class="has-text-centered">
        <loader/>
      </div>
    </template>
    <template v-else-if="$fetchState.error">
      <div class="notification is-danger">
        {{ $fetchState.error.message }}
      </div>
    </template>
    <template v-else-if="hundreds.length">
      <hundreds-list :hundreds="hundreds" visible-count="50" cards-in-column="3"/>
    </template>
    <template v-else-if="cards.length">
      <p class="has-text-centered">total: {{cards.length}}</p>
      <cardlist :cards="cards" cards-in-column="3" card-size="s" visible-count="33"/>
    </template>
  </section>
</template>

<script>
import HundredsList from "@/components/search/hundreds-list.vue";
import Cardlist from "@/components/list/cardlist.vue";

export default {
  components: {Cardlist, HundredsList},
  head() {
    let rank = parseInt(this.$route.params.rank)
    let hundred = parseInt(this.$route.params.hundred)
    return {
      title: `${this.$t('title')} ${hundred}-${hundred + rank - 1}`
    }
  },
  data() {
    return {
      hundreds: [],
      cards: []
    };
  },
  computed: {
    subtitle() {
      let rank = parseInt(this.$route.params.rank)
      let hundred = parseInt(this.$route.params.hundred)
      return `${this.$t('subtitle')} ${hundred}-${hundred + rank - 1}`
    }
  },
  async fetch() {
    let rank = parseInt(this.$route.params.rank)
    let hundred = parseInt(this.$route.params.hundred)
    if (rank === 100) {
      this.cards = await this.$axios.$get(`/search/${rank}/${hundred}`)
    } else {
      this.hundreds = await this.$axios.$get(`/search/${rank}/${hundred}`)
    }
  }
}
</script>

<style scoped>

</style>
