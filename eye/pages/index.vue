<template>
  <div>
    <short/>
    <promotion/>
    <heart/>
    <section>
      <a id="cards"></a>
      <last-artworks/>
    </section>
  </div>
</template>
<script>
import LastArtworks from '@/components/last-cards/last-artworks.vue'
import Short from "@/components/description/short.vue";
import Launch from "@/components/description/launch.vue";

export default {
  components: {Launch, Short, LastArtworks},
  layout: 'artchitect',
  head() {
    return {
      title: this.$t('index_title')
    }
  }
}
</script>
<style>

</style>
