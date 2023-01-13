<template>
  <div>
    <description/>
    <section class="green-section has-text-centered">
      <p>
        🚀 <b>Artchitect</b> launched 11th Jan 2023!
      </p>
      <p>
        Visit Artchitect's NFT-store:
        <a href="https://opensea.io/collection/artchitect" target="_blank">Artchitect Collection / Opensea</a>
      </p>
    </section>
    <section>
      <last-artworks/>
    </section>
  </div>
</template>
<script>
import LastArtworks from '@/components/last-cards/last-artworks.vue'
import Description from "@/components/description/description";

export default {
  components: {Description, LastArtworks},
  layout: 'artchitect',
  head: {
    title: 'Artchitect'
  }
}
</script>
<style>

</style>
