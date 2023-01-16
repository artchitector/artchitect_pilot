<template>
  <div>
    <section>
      {{ $t('selection_description') }}
    </section>
    <section>
      <h1 class="is-size-3 has-text-centered mb-5">{{ $t('selection') }}</h1>
      <div v-if="count" class="has-text-centered mb-6">total <b>{{ count }}</b></div>
      <div class="notification is-primary" v-if="$fetchState.pending">
        {{ $t('loading') }}...
      </div>
      <div class="notification is-primary" v-else-if="$fetchState.error">
        {{ $fetchState.error.message }}
      </div>
      <div class="columns" v-else v-for="line in lines">
        <div class="column has-text-centered" v-for="id in line">
          <a href="#" @click.prevent="onSelect(id)">
            <img :src="`/api/image/s/${id}`"/>
          </a>
        </div>
      </div>
    </section>
    <viewer ref="viewer"/>
  </div>
</template>
<script>
import Viewer from "@/components/viewer/viewer";

export default {
  components: {Viewer},
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
  },
  methods: {
    onSelect (id) {
      this.$refs.viewer.show(this.selection, id);
    }
  }
}
</script>
<style lang="scss">
.image-container {
  position: relative;

}
</style>
