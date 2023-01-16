<template>
  <div>
    <short/>
    <section>
      <a id="cards"></a>
      <last-artworks/>
    </section>
  </div>
</template>
<script>
import LastArtworks from '@/components/last-cards/last-artworks.vue'
import Description from "@/components/description/description";
import Description_ru from "~/components/description/description_ru.vue";
import Short from "@/components/description/short.vue";
import Launch from "@/components/description/launch.vue";

export default {
  components: {Launch, Short, Description_ru, Description, LastArtworks},
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
