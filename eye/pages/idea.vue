<template>
  <idea-frame/>
</template>
<script>
export default {
  head() {
    return {
      title: this.$t('idea_title')
    }
  },
  computed: {
    locale() {
      return this.$i18n.locale
    }
  }
}
</script>
