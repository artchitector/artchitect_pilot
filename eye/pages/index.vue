<template>
  <div>
    <section class="green-section has-text-centered">
      <p>
        🚀 <b>Artchitect</b> {{$t('launched')}}<br/>
        <a class="button is-success mt-2" href="#cards">{{ $t('to_cards')}}</a>
        <a class="button is-success mt-2" href="/pray">{{ $t('to_pray')}}</a><br/>
      </p>
    </section>
    <description_ru v-if="$i18n.locale === 'ru'"/>
    <description v-else/>
    <div class="box">
      <div class="has-text-centered">
        <a class="button is-success" href="/pray">{{ $t('to_pray')}}</a>
      </div>
    </div>
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

export default {
  components: {Description_ru, Description, LastArtworks},
  layout: 'artchitect',
  head: {
    title: 'Artchitect - Цифровой храм'
  }
}
</script>
<style>

</style>
