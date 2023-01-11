<template>
  <div>
    <description/>
    <section class="red-section has-text-centered">
      <p>
        Now <b>artchitect</b> is under massive development. Wait some time.
      </p>
      <p>
        The first artwork is planned for January 2023.
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
