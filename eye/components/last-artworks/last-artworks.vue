<template>
  <div>
    <div class="is-pulled-right">
      <a href="#" v-if="!jsonVisible" @click.prevent="jsonVisible = !jsonVisible">show json</a>
      <div v-else>
        <span @click.prevent="jsonVisible = !jsonVisible">hide json</span>
      </div>
    </div>

    <pre v-if="jsonVisible">{{this.artworks}}</pre>
    <h3 class="is-size-4 has-text-centered mb-4" v-if="count > 0">last {{count}} artworks</h3>
    <h3 class="is-size-4"v-else>last artworks</h3>
    <div class="notification is-primary" v-if="!artworks.length && $fetchState.pending">
      loading...
    </div>
    <div class="notification is-danger" v-else-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>

    <div v-else>
      <div class="has-text-centered mb-4">
          <nuxt-link v-for="page in pages" :to="page.url">
            {{page.caption}}
          </nuxt-link>
      </div>
      <div class="columns" v-for="line in lines">
        <div class="column" v-for="artwork in line">
          <artwork-view :artwork="artwork"/>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import ArtworkView from "~/components/last-artworks/artwork-view.vue";

export default {
  components: {ArtworkView},
  data() {
    return {
      artworks: [],
      error: null,
      updater: null,
      jsonVisible: false
    }
  },
  async fetch() {
    this.artworks = await this.$axios.$get('/last_paintings/100')
  },
  mounted() {
    this.updater = setInterval(() => {this.$fetch()}, 5000)
  },
  beforeDestroy() {
    clearInterval(this.updater)
  },
  computed: {
    pages() {
      if (!this.artworks.length) {
        return []
      }
      const lastId = this.artworks[0].ID
      const pages = []
      let from, to
      for (let i = 0; i < 5; i++) {
        from = lastId - (i * 100) - 1
        to = from - 100
        pages.push({
          "url": `/list?from=${from}&to=${to}`,
          "caption": `(page${i+1}:${from}-${to})`
        })
      }
      return pages
    },
    count() {
      return this.artworks.length
    },
    lines() {
      if (!this.artworks.length) {
        return [];
      }
      const lines = [
        []
      ]
      for (let i = 0; i < this.artworks.length; i++) {
        const artwork = this.artworks[i]
        let lastLine = lines.length - 1
        if (lines[lastLine].length < 3) {
          lines[lastLine].push(artwork)
        } else {
          lines.push([artwork])
        }
      }
      return lines
    }
  }
}
</script>
