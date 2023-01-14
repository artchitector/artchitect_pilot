<template>
  <div>
    <section class="green-section has-text-centered">
      <p>
        🚀 <b>Artchitect</b> launched 15th Jan 2023!<br/>
        <a class="button is-success mt-2" href="#cards">Go to cards</a>
      </p>
    </section>
    <description/>
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
