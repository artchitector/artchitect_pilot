<i18n>
{
  "en": {
    "title": "Artchitect - details"
  },
  "ru": {
    "title": "Artchitect - подробнее"
  }
}
</i18n>
<template>
  <endescription v-if="locale === 'en'"/>
  <rudescription v-else/>
</template>
<script>
import Rudescription from "@/components/idea/ru/rudescription.vue";
import Endescription from "@/components/idea/en/endescription.vue";

export default {
  components: {Endescription, Rudescription},
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
