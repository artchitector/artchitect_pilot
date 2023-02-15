<i18n>
{
  "en": {
    "page_title": "Artchitect - Log in",
    "title": "Log in"
  },
  "ru": {
    "page_title": "Artchitect - Вход",
    "title": "Вход"
  }
}
</i18n>
<template>
  <section class="content">
    <h3 class="is-size-5">{{ $t('title') }}</h3>
    <script async src="https://telegram.org/js/telegram-widget.js?21" data-telegram-login="ArtchitectBot" data-size="large" data-auth-url="https://artchitect.space/api/login" data-request-access="write"></script>
  </section>
</template>

<script>
export default {
  name: "login",
  head() {
    return {
      title: this.$t('title')
    }
  },
  data() {
    return {
      code: ""
    }
  },
  methods: {
    submit() {
      alert(`Submited ${this.code}`)
    }
  }
}
</script>

<style scoped>

</style>
