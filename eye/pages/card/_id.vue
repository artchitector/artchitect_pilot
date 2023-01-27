<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      loading...
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>
    <div v-else-if="card">
      <h1 class="is-size-2 has-text-centered">{{$t('card')}} #{{ card.ID }}</h1>
      <p><span class="tag is-primary is-light">version {{ card.Version }}</span></p>
      <p>{{ created }}</p>
      <p>seed = {{ card.Spell.Seed }}</p>
      <p>tags = <i>{{ card.Spell.Tags }}</i></p>
      <img :src="`/api/image/xf/${card.ID}`"/>
    </div>
  </section>

</template>
<script>
import moment from "moment"
export default {
  head () {
    return {
      title: `Artchitect - card #${this.$route.params.id}`
    }
  },
  data () {
    return {
      card: null
    }
  },
  computed: {
    created() {
      return moment(this.card.CreatedAt).format("YYYY MMM Do HH:mm:ss")
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
