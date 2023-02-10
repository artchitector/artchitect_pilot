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
<!--  <idea-frame v-if="locale === 'en'"/>-->
  <endescription v-if="locale === 'ru'"/>
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
