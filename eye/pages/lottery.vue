<template>
  <div>
    <section>
      {{ $t('lottery_description') }}
    </section>
    <section>
      <div class="notification is-primary" v-if="$fetchState.pending">
        loading...
      </div>
      <div class="notification is-danger" v-if="error">
        {{ error }}
      </div>
      <div class="notification is-danger" v-if="$fetchState.error">
        {{ $fetchState.error.message }}
      </div>
      <lottery v-for="lottery in lotteries" :lottery="lottery"/>
    </section>
  </div>
</template>
<script>
import Lottery from "@/components/lottery/lottery";
export default {
  head: {
    title: 'Artchitect - lottery'
  },
  components: {Lottery},
  data () {
    return {
      lotteries: null,
      error: null,
    }
  },
  async fetch () {
    try {
      this.lotteries = await this.$axios.$get('/lottery/10')
    } catch (e) {
      if (!!e.response && !!e.response.data && !!e.response.data.error) {
        this.error = e.response.data.error;
      }
    }
  }
}
</script>
