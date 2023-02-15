<i18n>
{
  "en": {
    "title": "Artchitect - liked",
    "subtitle": "you liked"
  },
  "ru": {
    "title": "Artchitect - нравится",
    "subtitle": "вам нравится"
  }
}
</i18n>
<template>
  <section>
    <h1 class="has-text-centered is-size-4">{{ $t('subtitle') }}</h1>
    <template v-if="$fetchState.pending">
      <loader/>
    </template>
    <template v-else-if="$fetchState.error">
      <div class="notification is-danger">{{$fetchState.error.message}}</div>
    </template>
    <template v-else>
      <cardlist :cards="liked" cards-in-column="5" card-size="m" visible-count="30"/>
    </template>
  </section>
</template>
<script>
import Cardlist from "@/components/list/cardlist.vue";

export default {
  name: "liked",
  components: {Cardlist},
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      liked: [],
    }
  },
  async fetch() {
    this.liked = await this.$axios.$get('/liked')
  }
}
</script>

<style scoped>

</style>
