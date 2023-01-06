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
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>

    <div v-else>
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
    this.artworks = await this.$axios.$get('/last_paintings/20')
  },
  mounted() {
    this.updater = setInterval(() => {this.$fetch()}, 5000)
  },
  beforeDestroy() {
    clearInterval(this.updater)
  },
  computed: {
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
