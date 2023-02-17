<i18n>
{
  "en": {
    "title": "Artchitect - search",
    "subtitle": "search"
  },
  "ru": {
    "title": "Artchitect - поиск",
    "subtitle": "поиск"
  }
}
</i18n>

<template>
  <section>
    <h3 class="has-text-centered is-size-4">{{ $t('subtitle') }}</h3>
  </section>
</template>

<script>
export default {
  name: "search",
  head() {
    return {
      title: this.$t('title')
    }
  }
}
</script>

<style scoped>

</style>
