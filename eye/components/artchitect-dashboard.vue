<template>
  <div>
    <a href="#" class="is-pulled-right" @click.prevent="reload()">reload {{ stateHash }}</a>
    <p v-if="$fetchState.pending">
      reloading...
    </p>
    <p v-else-if="$fetchState.error">
      <span class="has-text-danger">{{ $fetchState.error.message }}</span>
    </p>
    <div v-else>
      <dashboard-state :state="state.State"/>
    </div>
  </div>
</template>
<script>
import MD5 from 'crypto-js/md5'
import DashboardState from '@/components/dashboard/dashboard-state'

export default {
  components: { DashboardState },
  data: () => ({
    state: null,
    stateHash: null,
    updating: false
  }),
  mounted () {
    const self = this
    setInterval(async () => {
      // update state every 200ms.
      if (this.updating) {
        return
      }
      try {
        self.updating = true
        const newState = await this.$axios.$get('/state')
        console.log('loaded new state', newState)
        const newHash = MD5(JSON.stringify(newState)).toString()
        if (this.stateHash !== newHash) {
          console.log('old md5', this.stateHash, 'new md5', newHash, 'updating state...')
          self.state = newState
          self.stateHash = newHash
        } else {
          console.log('old md5', this.stateHash, 'new md5', newHash, 'not changed')
        }
      } finally {
        self.updating = false
      }
    }, 5000)
  },
  async fetch () {
    this.state = null
    this.state = await this.$axios.$get('/state')
    this.stateHash = MD5(JSON.stringify(this.state)).toString()
  },
  methods: {
    async reload () {
      await this.$fetch()
    }
  }
}
</script>
