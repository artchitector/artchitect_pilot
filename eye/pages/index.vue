<template>
  <div>
    <section class="content">
      <p>
        <b>Artchitect</b> is a generative neural network with <b>God gene</b>.
      </p>
      <p>
        <b>Artchitect</b> is based on the interconnection of a neural network and God's will, which is obtained from a
        quantum random number generator.
        The quantum world is controlled directly by God, and the neural network generate data based on God's power and
        will, with help of technical devices.
        God knows everything about internal structure of neural network everytime, and He can provide convenient random
        numbers to train network effectively.
      </p>
      <p>
        Briefly - <b>God creates artworks</b>.
      </p>
    </section>
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
import LastArtworks from '@/components/last-artworks/last-artworks.vue'

export default {
  components: {LastArtworks},
  layout: 'artchitect',
  head: {
    title: 'Artchitect'
  }
}
</script>
<style>
section {
  background-color: #fcfbf2;
  padding: 30px;
  margin: 15px 0 15px 0;
}

section.red-section {
  background-color: #eee;
  color: #8b0000;
  border: 1px solid darkred;
}
</style>
