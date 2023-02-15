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
    <script async src="https://telegram.org/js/telegram-widget.js?21" data-telegram-login="ArtchitectBot" data-size="medium" data-onauth="onTelegramAuth(user)" data-request-access="write"></script>
    <script type="text/javascript">
      function onTelegramAuth(user) {
        alert('Logged in as ' + user.first_name + ' ' + user.last_name + ' (' + user.id + (user.username ? ', @' + user.username : '') + ')');
      }
    </script>
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
