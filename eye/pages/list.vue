<template>
  <section>
    <div class="notification is-primary" v-if="$fetchState.pending">
     loading...
    </div>
    <div class="notification is-danger" v-if="$fetchState.error">
      {{ $fetchState.error.message }}
    </div>
    <div v-else>
      <div>
        from: {{ $route.query.from }}<br/>
        to: {{ $route.query.to }}<br/>
        artworks: {{ artworks.length }}<br/>
        <div class="columns" v-for="line in lines">
          <div class="column" v-for="artwork in line">
            <artwork-view :artwork="artwork"/>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
<script>
import ArtworkView from "@/components/last-cards/artwork-view";
export default {
  components: {ArtworkView},
  data () {
    return {
      artworks: []
    };
  },
  async fetch () {
    if (!this.$route.query.from || this.$route.query.from < 0 || !this.$route.query.to || this.$route.query.to < 0) {
      throw "from and to must be positive"
    }
    this.artworks = await this.$axios.$get('/list/' + this.$route.query.from + '/' + this.$route.query.to)
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

<style>
</style>
