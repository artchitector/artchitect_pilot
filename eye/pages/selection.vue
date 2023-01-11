<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
      loading...
    </div>
    <div class="notification is-primary" v-else-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>
    <div class="columns" v-else v-for="line in lines">
      <div class="column" v-for="id in line">
        <NuxtLink :to="`/card/${id}`" target="_blank">
          <img :src="`/api/painting/${id}`"/>
        </NuxtLink>
      </div>
    </div>
  </section>
</template>
<script>
export default {
  data() {
    return {
      selection: []
    }
  },
  async fetch () {
    this.selection = await this.$axios.$get("/selection")
  },
  computed: {
    lines() {
      const lines = []
      const chunkSize = 5
      for (let i = 0; i < this.selection.length; i += chunkSize) {
        lines.push(this.selection.slice(i, i + chunkSize))
      }
      return lines
    }
  }
}
</script>
