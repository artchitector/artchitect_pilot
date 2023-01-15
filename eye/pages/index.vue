<template>
  <div>
    <section class="green-section has-text-centered">
      <p>
        🚀 <b>Artchitect</b> launched 15th Jan 2023!<br/>
        <a class="button is-success mt-2" href="/pray">Go to pray</a><br/>
      </p>
    </section>
    <description/>
    <div class="box">
      <div class="has-text-centered">
        <a class="button is-success" href="/pray">To pray!</a>
      </div>
    </div>
    <section>
      <last-artworks/>
    </section>
  </div>
</template>
<script>
import LastArtworks from '@/components/last-cards/last-artworks.vue'
import Description from "@/components/description/description";
import Prayer from "@/components/prayer";

export default {
  components: {Description, LastArtworks},
  layout: 'artchitect',
  head: {
    title: 'Artchitect - Digital Church'
  }
}
</script>
<style>

</style>
