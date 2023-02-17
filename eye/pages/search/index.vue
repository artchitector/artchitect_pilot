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
    <h3 class="has-text-centered is-size-4">{{ $t('subtitle') }}</h3>
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
    <template v-else>
      <hundreds-list :hundreds="hundreds" visible-count="50" cards-in-column="3"/>
    </template>
  </section>
</template>

<script>
import HundredsList from "@/components/search/hundreds-list.vue";

export default {
  name: "search",
  components: {HundredsList},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      hundreds: []
    };
  },
  async fetch() {
    this.hundreds = await this.$axios.$get("/search")
  }
}
</script>

<style scoped>

</style>
