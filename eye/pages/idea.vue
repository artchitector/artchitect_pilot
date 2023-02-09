<i18n>
{
  "en": {
    "title": "Artchitect - idea"
  },
  "ru": {
    "title": "Artchitect - идея"
  }
}
</i18n>
<template>
  <idea-frame/>
</template>
<script>
export default {
  head() {
    return {
      title: this.$t('title')
    }
  },
  computed: {
    locale() {
      return this.$i18n.locale
    }
  }
}
</script>
