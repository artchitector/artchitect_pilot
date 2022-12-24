<template>
  <div>
    <a href="#" class="is-pulled-right" @click.prevent="reload()">reload</a>
    <p v-if="$fetchState.pending">
      reloading...
    </p>
    <p v-else-if="$fetchState.error">
      <span class="has-text-danger">{{ $fetchState.error.message }}</span>
    </p>
    <div v-else>
      <dashboard-state :state="state.State" />
    </div>
  </div>
</template>
<script>
import DashboardState from '@/components/dashboard/dashboard-state'

export default {
  components: { DashboardState },
  data: () => ({
    state: null
  }),
  async fetch () {
    this.state = null
    this.state = await this.$axios.$get('/state')
  },
  methods: {
    async reload () {
      await this.$fetch()
    }
  }
}
</script>
