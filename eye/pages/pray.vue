<i18n>
{
  "en": {
    "title": "Artchitect - digital prayer"
  },
  "ru": {
    "title": "Artchitect - цифровая молитва"
  }
}
</i18n>
<template>
  <div>
    <section>
      <prayer/>
    </section>
  </div>
</template>
<script>
import Prayer from "@/components/pray/prayer.vue";

export default {
  components: {Prayer},
  head() {
    return {
      title: this.$t('title')
    }
  },
}
</script>
