<template>
  <div>
    <section>
      {{ $t('selection_description') }}
    </section>
    <section>
      <h1 class="is-size-3 has-text-centered mb-5">{{$t('selection')}}</h1>
      <div v-if="count" class="has-text-centered mb-6">total <b>{{ count }}</b></div>
      <div class="notification is-primary" v-if="$fetchState.pending">
        {{$t('loading')}}...
      </div>
      <div class="notification is-primary" v-else-if="$fetchState.error">
        {{ $fetchState.error.message }}
      </div>
      <div class="columns" v-else v-for="line in lines">
        <div class="column" v-for="id in line">
          <NuxtLink :to="`/card/${id}`" target="_blank">
            <img :src="`/api/image/s/${id}`"/>
          </NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>
<script>
export default {
  head: {
    title: 'Artchitect - Выбор'
  },
  data () {
    return {
      selection: []
    }
  },
  async fetch () {
    this.selection = await this.$axios.$get("/selection")
  },
  computed: {
    lines () {
      const lines = []
      const chunkSize = 5
      for (let i = 0; i < this.selection.length; i += chunkSize) {
        lines.push(this.selection.slice(i, i + chunkSize))
      }
      return lines
    },
    count () {
      return this.selection.length
    }
  }
}
</script>
<style lang="scss">
.image-container {
  position: relative;

}
</style>
