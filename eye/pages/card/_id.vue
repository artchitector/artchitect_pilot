<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      loading...
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>
    <div v-else-if="card">
      <h1 class="is-size-2 has-text-centered">Card #{{ card.ID }}</h1>
      <p>created: {{ card.CreatedAt }}</p>
      <p>spell seed={{card.Spell.Seed}}, tags=<i>{{card.Spell.Tags}}</i></p>
      <img :src="`/api/image/f/${card.ID}`"/>
    </div>
  </section>

</template>
<script>
export default {
  head() {
    return {
      title: `Artchitect - card #${this.$route.params.id}`
    }
  },
  data () {
    return {
      card: null
    }
  },
  async fetch () {
    const id = parseInt(this.$route.params.id);
    if (!id) {
      throw "id must be positive integer"
    }
    this.card = await this.$axios.$get(`/card/${id}`)
  }
}
</script>
