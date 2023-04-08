<i18n>
{
  "en": {
    "title": "Artchitect - digital AI-artist"
  },
  "ru": {
    "title": "Artchitect - цифровой творец"
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
