<template>
  <div>
    <section class="green-section has-text-centered">
      <p>
        🚀 <b>Artchitect</b> launched 15th Jan 2023!<br/>
        <a class="button is-success mt-2" href="#cards">Go to cards</a><br/>
      </p>
    </section>
    <description/>
    <prayer/>
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
  components: {Description, LastArtworks, Prayer},
  layout: 'artchitect',
  head: {
    title: 'Artchitect'
  }
}
</script>
<style>

</style>
