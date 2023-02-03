<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      loading...
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>
    <div v-else-if="card">
      <h1 class="is-size-2 has-text-centered">{{ $t('card') }} #{{ card.ID }}</h1>
      <p><span class="tag is-primary is-light">version {{ card.Version }}</span></p>
      <p>{{ created }}</p>
      <p>seed = {{ card.Spell.Seed }}</p>
      <p class="tags">tags = <i>{{ card.Spell.Tags }}</i></p>
      <p class="has-text-centered">
        <a :href="fullSizeUrl" target="_blank" class="is-size-7">view full size</a>
      </p>
      <img :src="`/api/image/f/${card.ID}`"/>
    </div>
  </section>

</template>
<script>
import moment from "moment"

export default {
  head() {
    return {
      title: `Artchitect - card #${this.$route.params.id}`
    }
  },
  data() {
    return {
      card: null
    }
  },
  computed: {
    created() {
      return moment(this.card.CreatedAt).format("YYYY MMM Do HH:mm:ss")
    },
    fullSizeUrl() {
      return `${process.env.STORAGE_URL}/cards/card-${this.card.ID}.jpg`
    }
  },
  async fetch() {
    const id = parseInt(this.$route.params.id);
    if (!id) {
      throw "id must be positive integer"
    }
    this.card = await this.$axios.$get(`/card/${id}`)
  }
}
</script>
<style lang="scss" scoped>
p.tags {
  word-wrap: break-word;
  word-break: break-all;
  overflow: hidden;
}
</style>
