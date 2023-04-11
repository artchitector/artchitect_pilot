<i18n>
{
  "en": {
    "title": "Artchitect - dreaming AI"
  },
  "ru": {
    "title": "Artchitect - спящий ИИ, который видит сны"
  }
}
</i18n>
<template>
  <div>
    <short/>
<!--    <promotion/> temporary disabled-->
    <heart/>
    <section>
      <a id="cards"></a>
      <last99/>
    </section>
  </div>
</template>
<script>
import Short from "@/components/index/short.vue";
import Promotion from "@/components/index/promotion.vue";
import Last99 from "@/components/index/last99.vue";

export default {
  components: {Last99, Short, Promotion},
  layout: 'artchitect',
  head() {
    return {
      title: this.$t('title')
    }
  }
}
</script>
<style>

</style>
