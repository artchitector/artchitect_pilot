<template>
  <div>
    <section>
      {{ $t('lottery_description') }}
    </section>
    <section>
      <lottery-list/>
    </section>
  </div>
</template>
<script>
export default {
  head() {
    return {
      title: this.$t('lottery_title')
    }
  }
}
</script>
